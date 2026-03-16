package comment_service

import (
	"blogx_server/global"
	"blogx_server/models"

	"gorm.io/gorm"
)

func GetRootComment(commentID uint) (model *models.CommentModel) {
	var comment models.CommentModel
	err := global.Db.Take(&comment, commentID).Error
	if err != nil {
		return nil
	}

	if comment.ParentID == nil {
		// 没有父评论了，那他就是根评论
		return &comment
	}
	return GetRootComment(*comment.ParentID)
}

// GetCommentTree 获取评论树
func GetCommentTree(model *models.CommentModel) {
	global.Db.Preload("SubCommentList").Take(model)
	for _, commentModel := range model.SubCommentList {
		GetCommentTree(commentModel)
	}
}

// GetCommentTreeV2 获取评论树
func GetCommentTreeV2(id uint) (model *models.CommentModel) {
	model = &models.CommentModel{
		Model: models.Model{gorm.Model{ID: id}},
	}

	global.Db.Preload("SubCommentList").Take(model)
	for i := 0; i < len(model.SubCommentList); i++ {
		commentModel := model.SubCommentList[i]
		item := GetCommentTreeV2(commentModel.ID)
		model.SubCommentList[i] = item
	}
	return
}
