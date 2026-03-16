package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/comment_service"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()

	rootComment := comment_service.GetRootComment(1)
	fmt.Println(rootComment.ID)
	rootComment = comment_service.GetRootComment(2)
	fmt.Println(rootComment.ID)
	rootComment = comment_service.GetRootComment(3)
	fmt.Println(rootComment.ID)

}
