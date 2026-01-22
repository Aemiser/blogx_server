package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/log_service"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()
	flags.Run()

	log := log_service.NewRuntimeLog("test1", log_service.RuntimedateHour)
	log.SetItem("文章1", "test1")
	log.Save()
	log.SetItem("文章2", "test2")
	log.Save()
}
