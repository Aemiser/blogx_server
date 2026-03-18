package models

import "blogx_server/models/enum/message_type_enum"

type MessageModel struct {
	Model
	Type               message_type_enum.Type
	Message            string `json:"message"`
	RecvUserID         uint   `json:"revuser_id"`
	ActionUserID       uint   `json:"actionuser_id"`
	ActionUserNickName string `json:"actionuser_nickname"`
	ActionUserAvatar   string `json:"actionuser_avatar"`
	ArticleID          uint   `json:"article_id"`
	ArticleTitle       string `json:"article_title"`
	CommentID          uint   `json:"comment_id"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	LinkTitle          string `json:"linktitle"`
	LinkHref           string `json:"linkhref"`
	IsRead             bool   `json:"is_read"`
}
