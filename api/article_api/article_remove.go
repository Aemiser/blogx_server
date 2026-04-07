package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"blogx_server/service/message_service"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleRemoveView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDListRequest](c)

	log := log_service.GetLog(c)
	log.SetLogType(enum.OperationLogType)
	log.SetTitle("<span style='color: #ff4d4f'>🗑️ 管理员删除文章</span>")
	log.SetItem("操作类型", "<span style='color: #ff4d4f'>批量删除</span>")

	var list []models.ArticleModel
	global.Db.Find(&list, "id in  ?", cr.IDList)

	log.SetItem("请求删除数量", fmt.Sprintf("<span style='color: #faad14'>%d</span>", len(cr.IDList)))
	log.SetItem("实际查询到", fmt.Sprintf("<span style='color: #1890ff'>%d</span>", len(list)))

	if len(list) > 0 {
		var titles []string
		for _, model := range list {
			titles = append(titles, model.Title)
			message_service.InsertSystemMessage(model.UserID, "管理员删除了你的文章", fmt.Sprintf("%s 文章不符合社区规范", model.Title), "", "")
		}
		log.SetItem("文章标题列表", fmt.Sprintf("<div style='background: #f5f5f5; padding: 8px; border-radius: 4px; margin-top: 4px'>%s</div>", fmt.Sprintf("<br>%s", strings.Join(titles, "<br>"))))

		err := global.Db.Delete(&list).Error
		if err != nil {
			log.SetItemError("删除失败", err)
			res.FailWithCodeAndMsg(res.ArticleDeleteFail, "删除文章失败", c)
			return
		}
		log.SetItem("删除结果", fmt.Sprintf("<span style='color: #52c41a; font-weight: bold'>✅ 成功删除 %d 篇文章</span>", len(list)))
		log.SetItem("站内信通知", fmt.Sprintf("<span style='color: #52c41a'>✅ 已向 %d 位作者发送通知</span>", len(list)))
	} else {
		log.SetItem("删除结果", "<span style='color: #8c8c8c'>未找到任何文章</span>")
	}

	res.SuccessWithMsgf(c, "删除成功，成功删除文章 %d 条", len(list))
	return
}
