package models

import "gorm.io/gorm"

type Model struct {
	gorm.Model
}

type IDRequest struct {
	ID uint `json:"id" form:"id" uri:"id"`
}
