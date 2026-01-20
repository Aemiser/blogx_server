package models

type LogModel struct {
	Model
	LogType   int8      `json:"logType"` // 日志类型 1 2 3
	Title     string    `gorm:"size:64" json:"title"`
	Content   string    `json:"content"`
	Level     int8      `json:"level"` // 日志等级 1 2 3
	UserID    uint      `json:"userID"`
	UserModel UserModel `gorm:"foreignKey:UserID" json:"-"`
	IP        string    `gorm:"size:32" json:"ip"`
	Addr      string    `gorm:"size:64" json:"addr"`
	IsRead    bool      `json:"isRead"` //是否读取
}
