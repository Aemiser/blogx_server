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
	fmt.Println("评论的文章作者UID：", model.ArticleModel.UserID)
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
	global.Db.Preload("UserModel").Preload("ArticleModel").Take(&model)
	fmt.Println("里父评论ID：", model.ID)
	// 确定接收消息的用户ID
	var recvUserID uint

	// 如果有父评论，给父评论的作者发消息
	if model.ParentID != nil {
		var parentComment models.CommentModel
		err := global.Db.Preload("UserModel").Take(&parentComment, *model.ParentID).Error
		if err == nil {
			recvUserID = parentComment.UserID
			fmt.Println("给父评论作者发消息，父评论ID：", parentComment.ID, "作者UID：", recvUserID)
		} else {
			logrus.Warnf("查找父评论 %d 失败：%v", *model.ParentID, err)
			return
		}
	} else {
		// 如果没有父评论，说明是直接回复根评论，给根评论作者发消息
		recvUserID = model.UserID
		fmt.Println("直接回复根评论，给评论作者发消息，评论ID：", model.ID, "作者UID：", recvUserID)
	}

	err := global.Db.Create(&models.MessageModel{
		Type:               message_type_enum.CommentType,
		RecvUserID:         recvUserID,
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
