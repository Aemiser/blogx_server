package models

import "gorm.io/gorm"

type Model struct {
	gorm.Model
}

type IDRequest struct {
	ID uint `json:"id" form:"id" uri:"id"`
}
type IDListRequest struct {
	IDList []uint `json:"IDList" form:"IDList"`
}

type OptionsResponse[T any] struct {
	Label string `json:"label"`
	Value T      `json:"value"`
}
