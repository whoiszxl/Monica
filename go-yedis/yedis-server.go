package main

import (
	"Monica/go-yedis/utils"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	//默认配置文件路径
	defaultConfigPath = "yedis.conf"
)

//创建服务端实例
var zedis = new(core.Server)

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

	//读取命令行输入的ip和端口
	var netBind = flag.String("ip", netConfig.NetBind, "redis服务端IP")
	var netPort = flag.String("port", netConfig.NetPort, "redis服务端PORT")
	flag.Parse()
	host := *netBind + ":" + *netPort
	log.Println("Redis实例化地址：" + host)

	//监听退出事件做相应处理
	utils.ExitHandler()

	//初始化服务端实例
	initServer()

	//初始化网络监听并延时关闭
	netListener, err := net.Listen("tcp", host)
	if err != nil {
		log.Println("net listen err:", err)
	}
	defer netListener.Close()

	//循环监听新连接，将新连接放入go协程中处理
	for {
		conn, err := netListener.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

//处理连接请求
func handle(conn net.Conn) {
	c := yedis.CreateClient()
}

//初始化服务端实例
func initServer() {

}
