package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"blogx_server/service/message_service"
	"blogx_server/service/text_service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArticleExamineRequest struct {
	ArticleID uint   `json:"articleID" binding:"required"`
	Status    uint   `json:"status" binding:"required"`
	Msg       string `json:"msg"`
}

func (ArticleApi) ArticleExamineView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleExamineRequest](c)

	log := log_service.GetLog(c)
	log.SetLogType(enum.OperationLogType)
	log.SetTitle("<span style='color: #1890ff'>📝 文章审核</span>")
	log.SetItem("请求ID", fmt.Sprintf("<span style='color: #ff4d4f'>%d</span>", cr.ArticleID))

	var article models.ArticleModel
	err := global.Db.Take(&article, cr.ArticleID).Error
	if err != nil {
		res.FailWithCodeAndMsg(res.ArticleNotFound, "文章不存在", c)
		return
	}

	log.SetItem("文章标题", fmt.Sprintf("<span style='color: #52c41a'>%s</span>", article.Title))
	log.SetItem("文章作者ID", fmt.Sprintf("<span style='color: #722ed1'>%d</span>", article.UserID))
	log.SetItem("原状态", fmt.Sprintf("<span style='color: #faad14'>%d</span>", article.Status))

	statusText := "未知"
	statusColor := "#8c8c8c"
	switch cr.Status {
	case 3:
		statusText = "✅ 审核通过"
		statusColor = "#52c41a"
	case 4:
		statusText = "❌ 审核拒绝"
		statusColor = "#ff4d4f"
	}
	log.SetItem("新状态", fmt.Sprintf("<span style='color: %s; font-weight: bold'>%s</span>", statusColor, statusText))
	log.SetItem("审核备注", cr.Msg)

	// 修改状态
	global.Db.Model(&article).Update("status", cr.Status)
	switch cr.Status {
	case 3: // 审核成功
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", "文章审核成功", article.Title, fmt.Sprintf("/article/%d", article.ID))
		textList := text_service.MdContentTransformation(article.ID, article.Title, article.Content)

		err = global.Db.Create(&textList).Error
		if err != nil {
			logrus.Errorf("创建文章ID为：%d的text失败：%v", article.ID, err)
			log.SetItemError("创建全文索引失败", err)
		} else {
			log.SetItem("全文索引", "<span style='color: #52c41a'>✅ 创建成功</span>")
		}
	case 4: // 审核失败
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", fmt.Sprintf("文章审核失败，失败原因：%s", cr.Msg), "", "")
		log.SetItem("站内信通知", "<span style='color: #52c41a'>✅ 已发送</span>")
	}

	res.SuccessWithMsg("审核成功", c)

}
