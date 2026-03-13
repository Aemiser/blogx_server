package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleRemoveView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDListRequest](c)

	var list []models.ArticleModel
	global.Db.Find(&list, "id in  ?", cr.IDList)

	if len(list) > 0 {
		err := global.Db.Delete(&list).Error
		if err != nil {
			res.FailWithMsg("删除文章失败", c)
			return
		}
	}

	res.SuccessWithMsgf(c, "删除成功，成功删除文章 %d 条", len(list))
	return
}
