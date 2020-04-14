package utils

import (
	"fmt"
	"os"
)

//处理基础操作，查询版本，使用指南等
func BaseHelp(command string) {
	switch command {
	case "-v":
		version()
	case "-version":
		version()
	case "version":
		version()
	case "-h":
		help()
	case "-help":
		help()
	case "help":
		help()
	default:
		printError()
	}
}

func version() {
	fmt.Println("Redis server v=3.2.12 sha=00000000:0 malloc=jemalloc-3.6.0 bits=64 build=7897e7d0e13773f")
	os.Exit(0)
}

func help() {
	fmt.Println("Usage: ./yedis-server [/path/to/redis.conf] [options]")
	fmt.Println("       ./yedis-server - (read config from stdin)")
	fmt.Println("       ./yedis-server -v or --version")
	fmt.Println("       ./yedis-server -h or --help")
	fmt.Println("       ./redis-server --test-memory <megabytes>")
	fmt.Println("Examples:")
	fmt.Println("       ./yedis-server (run the server with default conf)")
	fmt.Println("       ./yedis-server /etc/redis/6380.conf")
	fmt.Println("       ./yedis-server --port 7777")
	fmt.Println("       ./yedis-server --port 7777 --slaveof 127.0.0.1 8888")
	fmt.Println("       ./yedis-server /etc/myredis.conf --loglevel verbose")
	fmt.Println("Sentinel mode:")
	fmt.Println("       ./yedis-server /etc/sentinel.conf --sentinel")
	os.Exit(0)
}

func printError() {
	fmt.Println("命令输入错误")
	fmt.Println(os.Getpid())
}
