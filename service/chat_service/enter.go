package chat_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/ctype/chat_type"
	"blogx_server/models/enum/chat_msg_type"
	"blogx_server/utils/xss"

	"github.com/sirupsen/logrus"
)

// ToChat A给B发消息
func ToChat(A, B uint, msgType chat_msg_type.MsgType, msg chat_type.ChatMsg) {
	chat := models.ChatModel{
		SeedUserID: A,
		RevUserID:  B,
		MsgType:    msgType,
		Msg:        msg,
	}

	err := global.Db.Create(&chat).Error
	if err != nil {
		logrus.Errorf("创建聊天记录失败: %v", err)
	}
}

func ToTextChat(A, B uint, content string) {
	ToChat(A, B, chat_msg_type.TextMsgType, chat_type.ChatMsg{
		ContentMsg: &chat_type.ContentMsg{
			Content: content,
		},
	})
}

func ToImageChat(A, B uint, src string) {
	ToChat(A, B, chat_msg_type.TextMsgType, chat_type.ChatMsg{
		ImagetMsg: &chat_type.ImagetMsg{
			Src: src,
		},
	})
}

func ToMarkdownChat(A, B uint, content string) {
	// 过滤xss
	filterContent := xss.Filter(content)
	ToChat(A, B, chat_msg_type.TextMsgType, chat_type.ChatMsg{
		MarkdownMsg: &chat_type.MarkdownMsg{
			Content: filterContent,
		},
	})
}
