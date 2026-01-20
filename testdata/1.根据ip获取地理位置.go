package main

import (
	"blogx_server/core"
	"fmt"
)

func main() {
	core.InItIPDB()
	fmt.Println(core.GetIpAddr("8.141.190.223"))
	fmt.Println(core.GetIpAddr("8.141.190.22"))
	fmt.Println(core.GetIpAddr("8.141.190.224"))
	fmt.Println(core.GetIpAddr("8.141.190.25"))
	fmt.Println(core.GetIpAddr("5.183.254.252"))
}
