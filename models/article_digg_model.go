package models

import (
	"time"

	"gorm.io/gorm"
)

// 用户对文章的点赞表
type ArticleDiggModel struct {
	Model
	UserID       uint           `gorm:"uniqueIndex:idx_name" json:"userID"`
	UserModel    UserModel      `gorm:"foreignKey:UserID" json:"-"`
	ArticleID    uint           `gorm:"uniqueIndex:idx_name" json:"articleID"`
	ArticleModel ArticleModel   `gorm:"foreignKey:ArticleID" json:"-"`
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt"`
}
