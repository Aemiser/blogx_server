package models

import (
	"blogx_server/models/ctype/chat_type"
	"blogx_server/models/enum/chat_msg_type"
)

type ChatModel struct {
	Model
	SeedUserID    uint                  `json:"seedUserID"`
	SeedUserModel UserModel             `gorm:"foreignKey:SeedUserID" json:"-"`
	RevUserID     uint                  `json:"RevUserID"`
	RevUserModel  UserModel             `gorm:"foreignKey:RevUserID" json:"-"`
	MsgType       chat_msg_type.MsgType `json:"msgType"` // 消息类型
	Msg           chat_type.ChatMsg     `gorm:"type:longtext;serializer:json" json:"msg"`
}
