package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

type ArticleExamineRequest struct {
	ArticleID uint   `json:"articleID" binding:"required"`
	Status    uint   `json:"status" binding:"required"`
	Msg       string `json:"msg"`
}

func (ArticleApi) ArticleExamineView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleExamineRequest](c)

	var article models.ArticleModel
	err := global.Db.Take(&article, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// 修改状态
	global.Db.Model(&article).Update("status", cr.Status)

	// TODO:给文章发布人发送一个系统消息

	res.SuccessWithMsg("审核成功", c)

}
