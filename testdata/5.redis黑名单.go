package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/redis_service/redis_jwt"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()
	global.Redis = core.InitRedis()

	//token, err := jwts.GetToken(jwts.Claims{
	//	UserID:   0,
	//	UserName: "taotao",
	//	Role:     0,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(token)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDbGFpbXMiOnsidXNlcklEIjowLCJ1c2VyTmFtZSI6InRhb3RhbyIsInJvbGUiOjB9LCJleHAiOjE3NjkyNzg1MzksImlzcyI6InRhb3RhbyJ9.bUXI3FKovMSrWRZC5ihcPWmSN-hoLx6ew8sxiLkvMJw"
	redis_jwt.BlackToken(token, redis_jwt.AdminBlackType)
	blk, ok := redis_jwt.HasTokenBlack(token)
	fmt.Println(blk, ok)

}
