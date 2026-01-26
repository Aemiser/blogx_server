package main

import (
	"blogx_server/utils/pwd"
	"fmt"
)

func main() {
	password := "123456"
	hashPassword, err := pwd.GenerateHashPassword(password)
	if err != nil {
		panic("生成加密密码失败")
	}
	fmt.Printf("生成的加密密码为:%s\n", hashPassword)
	if pwd.CompareHashAndPassword(hashPassword, password) {
		fmt.Println("匹配")
	} else {
		fmt.Println("不匹配")
	}
}
