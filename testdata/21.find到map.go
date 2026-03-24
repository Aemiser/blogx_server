package main

import (
	"blogx_server/common"
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"
)

func main() {
	// 初始化配置和数据库连接
	flags.Parse()
	global.Config = core.ReadConf()
	global.Db = core.InitDB()

	//mps := common.ScanMap(models.UserModel{}, common.ScanMapOptions{
	//	//	Where: global.Db.Where("id in ?", []uint{1, 2, 3}),
	//	//})
	mps := common.ScanMapV2(models.ChatModel{}, common.ScanMapOptions{
		Where: global.Db.Where("id in ?", []uint{1}),
	})
	fmt.Println(mps)
}
