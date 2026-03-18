package models

// 用户设定的文章分类表
type CategoryModel struct {
	Model
	Title       string         `gorm:"size:32" json:"title"`
	UserID      uint           `json:"userID"`
	UserModel   UserModel      `gorm:"foreignKey:UserID" json:"-"`
	ArticleList []ArticleModel `gorm:"foreignKey:CategoryID" json:"-"`
}
