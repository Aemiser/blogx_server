package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()

	//global.Db.Create(&models.ArticleModel{
	//	Title:   "测试文章",
	//	Content: "测试文章内容",
	//	TagList: ctype.List{"python", "go"},
	//})

	var list []models.ArticleModel
	global.Db.Find(&list)
	fmt.Println(list)

}
