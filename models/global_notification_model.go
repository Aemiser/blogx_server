package models

// 全局通知表
type GlobalNotificationModel struct {
	Model
	Title   string `gorm:"size:32" json:"title"`
	Content string `gorm:"size:63" json:"content"`
	Icon    string `gorm:"size:256" json:"icon"`
	Href    string `gorm:"size:256" json:"href"`
}
