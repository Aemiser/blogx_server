package message_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"
	"fmt"

	"github.com/sirupsen/logrus"
)

// InsertCommentMessage 文章评论消息
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
	fmt.Println("里父评论ID：", model.ID)
	fmt.Println("里2父评论ID：", model.ParentModel.UserID)
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

// InsertCollectArticleMessage 插入一条收藏文章消息
func InsertCollectArticleMessage(model models.UserArticleCollectModel) {
	global.Db.Preload("UserModel").Preload("ArticleModel").Take(&model)
	err := global.Db.Create(&models.MessageModel{
		Type:               message_type_enum.CollectArticleType,
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

// InsertDiggArticleMessage 插入一条评论点赞消息
func InsertDiggCommentMessage(model models.CommentDiggModel) {
	global.Db.Preload("CommentModel.ArticleModel").Preload("UserModel").Take(&model)
	err := global.Db.Create(&models.MessageModel{
		Type:               message_type_enum.DiggCommentType,
		RecvUserID:         model.CommentModel.UserID,
		ActionUserID:       model.UserID,
		ActionUserNickName: model.UserModel.Nickname,
		ActionUserAvatar:   model.UserModel.Avatar,
		ArticleID:          model.CommentModel.ArticleID,
		ArticleTitle:       model.CommentModel.ArticleModel.Title,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

func InsertSystemMessage(revUserID uint, title string, content string, linkTitle string, linkHref string) {
	err := global.Db.Create(&models.MessageModel{
		Type:       message_type_enum.SystemType,
		RecvUserID: revUserID,
		Title:      title,
		Content:    content,
		LinkTitle:  linkTitle,
		LinkHref:   linkHref,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}
