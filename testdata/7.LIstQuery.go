package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()

	var list []models.LogModel
	global.Db.Model(&models.LogModel{
		LogType: enum.LogInfoLevel,
	}).Find(&list)
	fmt.Println(list)
}
