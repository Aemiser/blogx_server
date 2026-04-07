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

type ArticleCreateRequest struct {
	Title       string             `json:"title" binding:"required"`
	Abstract    string             `json:"abstract"`
	Content     string             `json:"content" binding:"required"`
	CategoryID  *uint              `json:"categoryID"`
	TagList     ctype.List         `json:"tagList"`
	Cover       string             `json:"cover"`
	UserID      string             `json:"userID"`
	OpenComment bool               `json:"openComment"`
	Status      enum.ArticleStatus `json:"status" binding:"required,oneof=1 2"`
}

func (ArticleApi) ArticleCreateView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleCreateRequest](c)

	user, err := jwts.GetClaimsByGin(c).GetUser()
	if err != nil {
		res.FailWithCodeAndMsg(res.UserNotFound, "用户不存在", c)
		return
	}

	// 判断模式
	if global.Config.Site.SiteInfo.Mode == 2 {
		if user.Role != enum.AdminRole {
			res.FailWithMsg("博客模式下，用户无法发文章", c)
			return
		}
	}

	// 判断分类id是不是自己创建的
	var category models.CategoryModel
	if cr.CategoryID != nil {
		err = global.Db.Take(&category, "id  = ? and user_id = ?", cr.CategoryID, user.ID).Error
		if err != nil {
			res.FailWithCodeAndMsg(res.ArticleCategoryNotFound, "分类不存在", c)
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

	// 正文内容图片转存
	// 1.图片过多，同步做，接口耗时高，异步做，保存时间不确定，后续要对文章进行更新，写接口让前端做，没有保存完不能提交文章

	var article = models.ArticleModel{
		Title:       cr.Title,
		Abstract:    cr.Abstract,
		Content:     cr.Content,
		CategoryID:  cr.CategoryID,
		TagList:     cr.TagList,
		Cover:       cr.Cover,
		UserID:      user.ID,
		OpenComment: cr.OpenComment,
		Status:      cr.Status,
	}

	if cr.Status == 2 && global.Config.Site.Article.NoExamine {
		article.Status = enum.ArticlePublished
	}

	err = global.Db.Create(&article).Error
	if err != nil {
		res.FailWithCodeAndMsg(res.ArticleCreateFail, "文章创建失败", c)
		return
	}

	res.SuccessWithMsg("文章创建成功", c)

}
