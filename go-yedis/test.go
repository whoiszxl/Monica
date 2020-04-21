package main

import (
	"Monica/go-yedis/utils"
	"fmt"
)

func main() {

	fmt.Println(1 & 2)
	fmt.Println(2 & 2)
	fmt.Println(3 & 2)
	fmt.Println(1%50)
	fmt.Println(50%50)

	sec, ms := utils.CurrentSecondAndMillis()

	fmt.Println("时间戳", sec)
	fmt.Println("秒数", sec)
	fmt.Println("毫秒数", ms)
}
