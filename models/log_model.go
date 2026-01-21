package models

import "blogx_server/models/enum"

type LogModel struct {
	Model
	LogType     enum.LogType      `json:"logType"` // 日志类型 1 2 3
	Title       string            `gorm:"size:64" json:"title"`
	Content     string            `json:"content"`
	Level       enum.LogLevelType `json:"level"` // 日志等级 1 2 3
	UserID      uint              `json:"userID"`
	UserModel   UserModel         `gorm:"foreignKey:UserID" json:"-"`
	IP          string            `gorm:"size:32" json:"ip"`
	Addr        string            `gorm:"size:64" json:"addr"`
	IsRead      bool              `json:"isRead"`                  // 是否读取
	LoginStatus bool              `json:"loginStatus"`             // 登入日志的登录状态
	Username    string            `gorm:"size:32" json:"username"` // 登入日志的用户名
	Pwd         string            `gorm:"size:32" json:"pwd"`      // 登入日志的密码
	LoginType   enum.LoginType    `json:"loginType"`               // 登入的类型
}
