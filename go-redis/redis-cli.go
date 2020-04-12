package main

import "fmt"
import "flag"
import "os"
import "bufio"
import "net"
import "strings"
import "Monica/go-redis/error"

//获取客户端cmd命令行输入的ip和端口
var ip = flag.String("ip", "localhost", "redis服务端IP")
var port = flag.String("port", "6380", "redis服务端PORT")


func main() {
	//获取host
	flag.Parse()
	host := *ip + ":" + *port
	commandStart := host + "> "

	reader := bufio.NewReader(os.Stdin)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	error.CheckError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	error.CheckError(err)

	defer conn.Close()

	
	for {
		fmt.Print(commandStart)
		inputCmd, _ := reader.ReadString('\n')
		
		inputCmd = strings.Replace(inputCmd, "\n", "", -1)

		//TODO 将获取的命令发送到服务端
		//sendCmd2Server(text, conn)

		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		error.CheckError(err)

		if n == 0 {
			fmt.Println(commandStart, "nil")
		}else {
			fmt.Println(commandStart, string(buff))
		}
	}

}


