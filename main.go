package main

import (
	"blogx_server/common/res"
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/router"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()
	global.Redis = core.InitRedis()
	global.ESClient = core.EsConnect()
	flags.Run()

	core.InitMysqlEs()

	res.InitSysCode()
	res.InitUserCode()
	res.InitChatCode()
	res.InitImageCode()
	res.InitArticleCode()
	res.InitFocusCode()
	res.InitCommentCode()

	router.Run()
}
