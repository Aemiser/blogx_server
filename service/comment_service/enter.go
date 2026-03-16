package comment_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"time"

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

// GetParents 获取一个根评论的所有父评论
func GetParents(commentID uint) (list []models.CommentModel) {
	var comment models.CommentModel
	err := global.Db.Take(&comment, commentID).Error
	if err != nil {
		return nil
	}
	list = append(list, comment)
	if comment.ParentID != nil {
		// 还有父评论了 添加到列表中
		list = append(list, GetParents(*comment.ParentID)...)
	}
	return
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

func GetCommentOneDimensionalization(id uint) (list []models.CommentModel) {
	model := models.CommentModel{
		Model: models.Model{gorm.Model{ID: id}},
	}

	global.Db.Preload("SubCommentList").Take(&model)
	list = append(list, model)
	for _, item := range model.SubCommentList {
		subList := GetCommentOneDimensionalization(item.ID)
		list = append(list, subList...)
	}
	return
}

type CommentResponse struct {
	models.CommentModel
	ID           uint               `json:"id"`
	CreatedAt    time.Time          `json:"createdAt"`
	Content      string             `json:"content"`
	UserID       uint               `json:"userID"`
	UserNickname string             `json:"nickName"`
	UserAvatar   string             `json:"userAvatar"`
	ArticleID    uint               `json:"articleID"`
	ParentID     *uint              `json:"parentID"`
	DiggCount    int                `json:"diggCount"`
	ApplyCount   int                `json:"applyCount"`
	SubComments  []*CommentResponse `json:"subComments"`
}

func GetCommentTreeV4(id uint) (res *CommentResponse) {
	model := models.CommentModel{
		Model: models.Model{gorm.Model{ID: id}},
	}

	global.Db.Preload("UserModel").Preload("SubCommentList").Take(&model)

	res = &CommentResponse{
		ID:           model.ID,
		CreatedAt:    model.CreatedAt,
		Content:      model.Content,
		UserID:       model.UserID,
		UserNickname: model.UserModel.Nickname,
		UserAvatar:   model.UserModel.Avatar,
		ArticleID:    model.ArticleID,
		ParentID:     model.ParentID,
		DiggCount:    model.DiggCount,
		ApplyCount:   0,
	}
	for _, commentModel := range model.SubCommentList {
		res.SubComments = append(res.SubComments, GetCommentTreeV4(commentModel.ID))
	}
	return
}
