package site_msg_api

import (
	"blogx_server/common"
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
	var chatList []models.ChatModel

	// 如果接收人是我，而且这个消息未读
	global.Db.Find(&chatList, "rev_user_id = ?", userID) // 先找到接收人是我的所有消息
	var chatIDList []uint
	for _, model := range chatList {
		chatIDList = append(chatIDList, model.ID) // 获取id列表
	}

	chatAcMap := common.ScanMapV2(models.UserChatAtionModel{}, common.ScanMapOptions{
		Where: global.Db.Where("chat_id in ? ", chatIDList), // 在用户操作消息表中，根据id获取消息，并装进map中，key为id,value为model
	})
	for _, model := range chatList { // 遍历我所有的消息，不在map中的即为未读的，因为只有已读了的消息才会在用户操作消息表中插入消息
		_, ok := chatAcMap[model.ID]
		if !ok {
			data.PrivateMsgCount++
			continue
		}
	}

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
