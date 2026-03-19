package site_msg_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"

	"github.com/gin-gonic/gin"
)

type UserMsgResponse struct {
	CommentMsgCount int `json:"commentMsgCount"`
	DiggMsgCount    int `json:"diggMsgCount"`
	SystemMsgCount  int `json:"systemMsgCount"`
	PrivateMsgCount int `json:"privateMsgCount"`
}

func (SiteMsgApi) UserMsgView(c *gin.Context) {
	userID := jwts.GetUserIDByGin(c)

	var msgList []models.MessageModel
	global.Db.Find(&msgList, "recv_user_id  = ? and is_read = ? ", userID, false)

	var data UserMsgResponse

	for _, model := range msgList {
		switch model.Type {
		case message_type_enum.CommentType, message_type_enum.ApplyType:
			data.CommentMsgCount++
		case message_type_enum.DiggCommentType, message_type_enum.DiggArticleType, message_type_enum.CollectArticleType:
			data.DiggMsgCount++
		case message_type_enum.SystemType:
			data.SystemMsgCount++
		}
	}

	//TODO:未读私信的数量

	//算未读的全局消息
	// 过滤掉已读和删除的
	var userReadMsgList []uint
	global.Db.Model(&models.UserGlobalnotificationModel{}).
		Where("user_id =? and (is_read = ? or is_delete = ?)", userID, true, true).
		Select("id").Scan(&userReadMsgList)

	// 算未读的全局消息
	var systemMsg []models.GlobalNotificationModel
	query := global.Db.Where("")
	if len(userReadMsgList) > 0 {
		query.Where("id not in ? ", userReadMsgList)
	}
	global.Db.Where(query).Find(&systemMsg)
	data.SystemMsgCount += len(systemMsg)
	res.SuccessWithData(data, c)
}
