package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/service/message_service"
	"fmt"

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
	switch cr.Status {
	case 3: // 审核成功
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", "文章审核成功", article.Title, fmt.Sprintf("/article/%d", article.ID))
	case 4: // 审核失败
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", fmt.Sprintf("文章审核失败，失败原因：%s", cr.Msg), "", "")

	}

	res.SuccessWithMsg("审核成功", c)

}
