package main

import (
	"blogx_server/common/jwts"
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models/enum"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()
	global.Redis = core.InitRedis()
	token, err := jwts.GetToken(jwts.Claims{
		UserID:   0,
		UserName: "taotao",
		Role:     enum.AdminRole,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	claims, err := jwts.ParseToken(token)
	fmt.Println(claims, err)
	//claims, err := jwts.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDbGFpbXMiOnsidXNlcklEIjozLCJ1c2VyTmFtZSI6ImJfMDUzOSIsInJvbGUiOjF9LCJleHAiOjE3Njk4NTkxODAsImlzcyI6InRhb3RhbyJ9.CwPoxZZgT8MIW95cpBj8ZgJh1yFDB4-Bwc_tNJedZrs")
	//fmt.Println(claims, err)

}
