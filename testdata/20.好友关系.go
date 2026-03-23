package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/chat_service"
)

func main() {
	// 初始化配置和数据库连接
	flags.Parse()
	global.Config = core.ReadConf()
	global.Db = core.InitDB()

	// 执行测试
	//chat_service.ToTextChat(1, 2, "你好")
	chat_service.ToTextChat(1, 4, "你好")
	chat_service.ToImageChat(1, 4, "http://baidu.com")
	chat_service.ToImageChat(4, 1, "http://baidu.22")
}
