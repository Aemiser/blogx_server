package models

type UserGlobalnotificationModel struct {
	Model
	NotificationID uint `json:"notificationID"`
	UserID         uint `json:"userID"`
	IsRead         bool `json:"isRead"`
	IsDelete       bool `json:"isDelete"`
}
