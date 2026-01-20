package models

import "time"

type UserArticleCollectModel struct {
	Model
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userID"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"articleID"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
	CollectID    uint         `gorm:"uniqueIndex:idx_name" json:"collectID"` // 收藏夹ID
	CollectModel CollectModel `gorm:"foreignKey:CollectID" json:"-"`         // 属于哪个收藏夹
	CreatedAt    time.Time    `json:"createdAt"`
}
