package article_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

type ArticleRecommentResponse struct {
	ID        uint   `json:"id" gorm:"column:id"`
	Title     string `json:"title" gorm:"column:title"`
	LookCount int    `json:"lookCount" gorm:"column:look_count"`
}

func (ArticleApi) ArticleRecommentView(c *gin.Context) {
	cr := middlerware.GetBind[common.PageInfo](c)

	var list = make([]ArticleRecommentResponse, 0)
	global.Db.Model(models.ArticleModel{}).
		Select("id", "title", "look_count").
		Order("look_count desc").
		Limit(cr.GetLimit()).
		Scan(&list)

	res.SuccessWithList(list, len(list), c)
}
