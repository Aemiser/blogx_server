package main

import (
	"blogx_server/utils"
	"fmt"
)

func main() {
	var rlist = []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	list := utils.Reverse(rlist)
	fmt.Println(list)

}
