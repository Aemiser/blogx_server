package article_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	"blogx_server/utils"

	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleTagListView(c *gin.Context) {
	claims := jwts.GetClaimsByGin(c)

	var articleList []models.ArticleModel
	global.Db.Find(&articleList, "user_id = ? and status = ?", claims.Claims.UserID, enum.ArticlePublished)

	var tagList ctype.List
	for _, article := range articleList {
		tagList = append(tagList, article.TagList...)
	}

	tagList = utils.Unique(tagList)
	var list = make([]models.OptionsResponse[string], 0)
	for _, s := range tagList {
		list = append(list, models.OptionsResponse[string]{
			Label: s,
			Value: s,
		})
	}
	res.SuccessWithData(list, c)
}
