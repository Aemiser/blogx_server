package main

import (
	"blogx_server/models"
	"blogx_server/service/text_service"
	"fmt"
	"os"

	"gorm.io/gorm"
)

func main() {
	byteData, err := os.ReadFile("testdata/md/text.md")
	if err != nil {
		panic(err)
	}
	//MdContentTransformation(1, "这是一个大大的测试内容", string(byteData))
	list := text_service.MdContentTransformation(models.ArticleModel{
		Model:   models.Model{gorm.Model{ID: 1}},
		Title:   "这是一个大大的测试内容",
		Content: string(byteData),
	})

	fmt.Println(list)

}
