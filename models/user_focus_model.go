package models

type UserFocusModel struct {
	Model
	UserID         uint      `json:"userID"`
	UserModel      UserModel `gorm:"foreignKey:UserID" json:"-"`
	FocusUserID    uint      `json:"focusID"`
	FocusUserModel UserModel `gorm:"foreignKey:FocusUserID" json:"-"`
}
