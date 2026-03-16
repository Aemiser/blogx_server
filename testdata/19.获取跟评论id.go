package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"

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

	list := GetCommentOneDimensionalization(1)
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
