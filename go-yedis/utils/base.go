package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// 监听退出命令行事件
func ExitHandler() {

	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go signalHandler(c)
}

func signalHandler(c chan os.Signal) {
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			ExitFunc()
		default:
			fmt.Println("other", s)
		}
	}
}

func ExitFunc() {
	fmt.Println("Exit Yedis...")
	os.Exit(0)
}
