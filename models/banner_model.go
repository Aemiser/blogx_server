package models

type BannerModel struct {
	Model
	Cover string `gorm:"size:256" json:"cover"`
	Href  string `gorm:"size:256" json:"href"`
}
