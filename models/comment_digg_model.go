package models

import (
	"time"

	"gorm.io/gorm"
)

type CommentDiggModel struct {
	UserID       uint           `gorm:"uniqueIndex:idx_name" json:"userID"`
	UserModel    UserModel      `gorm:"foreignKey:UserID" json:"-"`
	CommentID    uint           `gorm:"uniqueIndex:idx_name" json:"commentID"`
	CommentModel CommentModel   `gorm:"foreignKey:CommentID" json:"-"`
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt"`
}
