package main

import (
	"blogx_server/common/jwts"
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"fmt"
)

func main() {
	// 启动服务
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	token, err := jwts.GetToken(jwts.Claims{
		UserID:   1,
		UserName: "admin",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("token:", token)
	//token := " eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDbGFpbXMiOnsidXNlcklEIjoxLCJ1c2VyTmFtZSI6ImFkbWluIiwicm9sZSI6MH0sImV4cCI6MTc2OTI1NDc2NiwiaXNzIjoidGFvdGFvIn0.x7nYnvQ9H9qMlRQETXou9nRNDiBxoUs7JcU_-Tqiz"
	claim, err := jwts.ParseToken("xx")
	fmt.Println(claim, err)

}
