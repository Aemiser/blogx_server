package comment_service

import (
	"blogx_server/global"
	"blogx_server/models"
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
