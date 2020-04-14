package main

import "fmt"
import "Monica/go-yedis/utils"

func main() {

	//获取配置
	netConfig, dbConfig, aofConfig := utils.ReadConfig()
	fmt.Println(netConfig.NetBind)
	fmt.Println(dbConfig)
	fmt.Println(aofConfig)
}