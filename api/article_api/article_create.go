package article_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum"

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
		res.FailWithMsg("用户不存在", c)
		return
	}

	// 判断分类id是不是自己创建的
	// 防止文章正文xss注入
	// 正文内容图片转存
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

	if global.Config.Site.Article.NoExamine {
		article.Status = enum.ArticlePublished
	}

	err = global.Db.Create(&article).Error
	if err != nil {
		res.FailWithMsg("文章创建失败", c)
		return
	}

	res.SuccessWithMsg("文章创建成功", c)

}
