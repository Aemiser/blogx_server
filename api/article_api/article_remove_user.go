package article_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleRemoveUserView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDRequest](c)

	claims := jwts.GetClaimsByGin(c)
	var article models.ArticleModel
	err := global.Db.Take(&article, "user_id = ? and id = ?", claims.Claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	err = global.Db.Delete(&article).Error
	if err != nil {
		res.FailWithMsg("删除文章失败", c)
		return
	}
	res.SuccessWithMsg("删除文章成功", c)
	return
}
