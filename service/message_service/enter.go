package message_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"

	"github.com/sirupsen/logrus"
)

func InsertCommentMessage(model models.CommentModel) {
	global.Db.Preload("UserModel").Preload("ArticleModel").Take(&model)
	err := global.Db.Create(&models.MessageModel{
		Type:               message_type_enum.CommentType,
		RecvUserID:         model.ArticleModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickName: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		ArticleID:          model.ArticleID,
		ArticleTitle:       model.ArticleModel.Title,
		CommentID:          model.ID,
		Content:            model.Content,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertApplyMessage 插入一条回复消息
func InsertApplyMessage(model models.CommentModel) {
	global.Db.Preload("ParentModel").Preload("UserModel").Preload("ArticleModel").Take(&model)
	err := global.Db.Create(&models.MessageModel{
		Type:               message_type_enum.CommentType,
		RecvUserID:         model.ParentModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickName: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		ArticleID:          model.ArticleID,
		ArticleTitle:       model.ArticleModel.Title,
		CommentID:          model.ID,
		Content:            model.Content,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertDiggArticleMessage 插入一条点赞文章消息
func InsertDiggArticleMessage(model models.ArticleDiggModel) {
	global.Db.Preload("UserModel").Preload("ArticleModel").Take(&model)
	err := global.Db.Create(&models.MessageModel{
		Type:               message_type_enum.DiggArticleType,
		RecvUserID:         model.ArticleModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickName: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		ArticleID:          model.ArticleID,
		ArticleTitle:       model.ArticleModel.Title,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}
