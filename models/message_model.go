package models

type MessageModel struct{
	ID uint `json:"ID"`
	Title string `json:"title"`
	Message string `json:"message"`
	RevUserID uint `json:"revuser_id"`
	ActionUserID uint`json:"actionuser_id"`
	ActionUserAvatar string `json:"actionuser_avatar"`
	LinkTitle string `linktitle`
	LinkHref string `linkhref`
}