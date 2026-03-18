package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()

	//model := models.CommentModel{
	//	Model: models.Model{gorm.Model{ID: 1}},
	//}
	//GetCommentTree(&model)
	//model := GetCommentTreeV3(1)
	//for _, c1 := range model.SubCommentList {
	//	fmt.Println("  ", c1.ID)
	//	for _, c2 := range c1.SubCommentList {
	//		fmt.Println("    ", c2.ID)
	//		for _, c3 := range c2.SubCommentList {
	//			fmt.Println("      ", c3.ID)
	//			for _, c4 := range c3.SubCommentList {
	//				fmt.Println("        ", c4.ID)
	//			}
	//		}
	//
	//	}
	//}

	//list := GetCommentOneDimensionalization(1)
	//for _, item := range list {
	//	fmt.Println(item.ID)
	//}
	//
	//res := GetCommentTreeV4(1)
	//byteData, _ := json.Marshal(res)
	//fmt.Println(string(byteData))

	list := GetParents(9)
	for _, item := range list {
		fmt.Println(item.ID)
	}
}

func GetCommentTree(model *models.CommentModel) {
	global.Db.Preload("SubCommentList").Take(model)
	for _, commentModel := range model.SubCommentList {
		GetCommentTree(commentModel)
	}
}

func GetCommentTreeV3(id uint) (model *models.CommentModel) {
	model = &models.CommentModel{
		Model: models.Model{gorm.Model{ID: id}},
	}

	global.Db.Preload("SubCommentList").Take(model)
	for i := 0; i < len(model.SubCommentList); i++ {
		commentModel := model.SubCommentList[i]
		item := GetCommentTreeV3(commentModel.ID)
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
