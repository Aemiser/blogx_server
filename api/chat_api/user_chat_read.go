package chat_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/ctype/chat_type"
	"blogx_server/models/enum/chat_msg_type"

	"github.com/gin-gonic/gin"
)

func (ChatApi) UserChatReadView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDRequest](c)
	userID := jwts.GetUserIDByGin(c)
	var chat models.ChatModel
	err := global.Db.Take(&chat, cr.ID).Error
	if err != nil {
		res.FailWithMsg("聊天记录不存在", c)
		return
	}
	item := ChatResponse{
		ChatListResponse: ChatListResponse{
			ChatModel: models.ChatModel{
				MsgType: chat_msg_type.MsgReadMsg,
				Msg: chat_type.ChatMsg{
					ReadMsg: &chat_type.ReadMsg{
						ChatID: chat.ID,
					},
				},
			},
		},
	}
	var chatAc models.UserChatAtionModel
	err = global.Db.Take(&chatAc, "user_id =? and chat_id = ?", userID, cr.ID).Error
	if err != nil { // 不存在即创建
		global.Db.Create(&models.UserChatAtionModel{
			UserID: userID,
			ChatID: cr.ID,
			IsRead: true,
		})
	}

	if chatAc.IsRead {
		res.SuccessWithMsg("消息读取成功", c)
		res.SendWsMsg(OnlineMap, chat.SeedUserID, item)
		return
	}

	if chatAc.IsDelete {
		res.FailWithMsg("消息已删除", c)
		return
	}

	global.Db.Model(&models.UserChatAtionModel{}).Where("user_id = ? and chat_id = ?", userID, chat.ID).Update("is_read", true)
	res.SendWsMsg(OnlineMap, chat.SeedUserID, item)
	res.SuccessWithMsg("消息读取成功", c)

}
