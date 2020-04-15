package main

import (
	"Monica/go-yedis/encrypt"
	"Monica/go-yedis/utils"
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	//获取客户端cmd命令行输入的IP和端口，默认值取: localhost:6380
	var ip = flag.String("ip", "localhost", "redis服务端IP")
	var port = flag.String("port", "6380", "redis服务端PORT")
	flag.Parse()
	host := *ip + ":" + *port
	commandStart := host + "> "

	//解析地址和端口号，创建一个TCPAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	utils.ErrorVerify("Tcp Addr创建失败", err, true)

	//建立连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	utils.ErrorVerify("Tcp 连接建立失败", err, true)

	defer conn.Close()

	for {
		//获取键盘输入
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(commandStart)
		cmd, err := reader.ReadString('\n')
		utils.ErrorVerify("命令读取失败", err, false)
		cmd = strings.Replace(cmd, "\n", "", -1)
		if cmd != "" {
			send2Server(cmd, conn)
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			resp, err := encrypt.DecodeFromBytes(buff)
			utils.ErrorVerify("服务端返回消息解析失败", err, false)

			if n == 0 {
				fmt.Println("nil")
			}else if err == nil {
				fmt.Println(string(resp.Value))
			}else {
				fmt.Println("服务器发生错误", err)
			}
		}
	}
}

//发送cmd命令到服务器端
func send2Server(cmd string, conn net.Conn) int{
	//给cmd转换为yedis协议格式
	cmdBytes, err := encrypt.EncodeCmd(cmd)
	utils.ErrorVerify("发送消息时编码失败", err, false)
	n, err := conn.Write(cmdBytes)
	utils.ErrorVerify("发送消息体失败", err, false)
	return n
}
