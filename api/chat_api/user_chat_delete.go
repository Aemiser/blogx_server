package chat_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func (ChatApi) UserChatDeleteView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDRequest](c)
	userID := jwts.GetUserIDByGin(c)

	var chat models.ChatModel
	err := global.Db.Take(&chat, cr.ID).Error
	if err != nil {
		res.FailWithMsg("消息不存在", c)
		return
	}

	// 之前是否操作过
	var chatAc models.UserChatAtionModel
	err = global.Db.Take(&chatAc, "user_id = ? and chat_id = ?", userID, chat.ID).Error
	// 不存在，则直接创建操作记录，如果是我自己删除的话，就不用管是不是已读
	if err != nil {
		global.Db.Create(&models.UserChatAtionModel{
			UserID:   userID,
			ChatID:   chat.ID,
			IsDelete: true,
		})
		res.SuccessWithMsg("消息删除成功", c)
		return
	}

	if chatAc.IsDelete {
		// 说明之前已经删除了
		res.SuccessWithMsg("消息已删除", c)
		return
	}

	global.Db.Model(&chatAc).Update("is_delete", true)
	res.SuccessWithMsg("消息删除成功", c)
	return

}
