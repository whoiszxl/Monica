package main

import (
	"Monica/go-yedis/utils"
	"fmt"
	"os"
	"strings"
)

const (
	//默认配置文件路径
	defaultConfigPath = "yedis.conf"
)

func main() {

	//获取用户输入的参数
	cmdArgs := os.Args

	//第一个参数是文件路径， 读取配置
	var configPath = defaultConfigPath
	if len(cmdArgs) > 1 && strings.LastIndex(cmdArgs[1], ".conf") != -1 {
		fmt.Println("读取配置")
		configPath = cmdArgs[1]
	} else if len(cmdArgs) == 2 {
		utils.BaseHelp(cmdArgs[1])
	}

	//获取配置
	netConfig, dbConfig, aofConfig := utils.ReadConfig(configPath)
	fmt.Println(netConfig)
	fmt.Println(dbConfig)
	fmt.Println(aofConfig)
}
