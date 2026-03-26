package article_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	"blogx_server/utils/markdown"
	"blogx_server/utils/xss"

	"github.com/gin-gonic/gin"
)

type ArticleUpdateRequest struct {
	ID          uint       `json:"ID" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Abstract    string     `json:"abstract"`
	Content     string     `json:"content" binding:"required"`
	CategoryID  *uint      `json:"categoryID"`
	TagList     ctype.List `json:"tagList"`
	Cover       string     `json:"cover"`
	UserID      string     `json:"userID"`
	OpenComment bool       `json:"openComment"`
}

func (ArticleApi) ArticleUpdateView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleUpdateRequest](c)

	user, err := jwts.GetClaimsByGin(c).GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// 判断模式
	if global.Config.Site.SiteInfo.Mode == 2 {
		if user.Role != enum.AdminRole {
			res.FailWithMsg("博客模式下，用户无法更新文章", c)
			return
		}
	}

	// 查找文章
	var article models.ArticleModel
	err = global.Db.Take(&article, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// 更新的文章必须是自己的
	if article.UserID != user.ID {
		res.FailWithMsg("更新的文章必须是自己的", c)
		return
	}
	// 判断分类id是不是自己创建的
	var category models.CategoryModel
	if cr.CategoryID != nil {
		err = global.Db.Take(&category, "id  = ? and user_id = ?", cr.CategoryID, user.ID).Error
		if err != nil {
			res.FailWithMsg("分类不存在", c)
			return
		}
	}

	// 防止文章正文xss注入
	cr.Content = xss.Filter(cr.Content)

	// 如果不传简介，从正文中取前30个字符
	if cr.Abstract == "" {
		cr.Content, err = markdown.ExtractContent(cr.Content, 200)
		if err != nil {
			res.FailWithMsg("文章正文解析失败", c)
			return
		}
	}

	mps := map[string]interface{}{
		"title":       cr.Title,
		"abstract":    cr.Abstract,
		"content":     cr.Content,
		"categoryID":  cr.CategoryID,
		"tagList":     cr.TagList,
		"cover":       cr.Cover,
		"openComment": cr.OpenComment,
	}
	if article.Status == enum.ArticlePublished && !global.Config.Site.Article.NoExamine {
		// 如果是已发布的文章进行编辑，则需要进行审核，如果是草稿则不变，如果是免审核状态，条件不达标也不修改
		mps["status"] = enum.ArticleExamine
	}
	// 更新文章
	err = global.Db.Model(article).Updates(mps).Error
	if err != nil {
		res.FailWithMsg("文章更新失败", c)
		return
	}
	res.SuccessWithMsg("文章更新成功", c)

}
