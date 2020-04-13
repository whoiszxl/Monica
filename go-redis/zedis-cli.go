package main

import "fmt"
import "flag"
import "os"
import "bufio"
import "net"
import "log"
import "strings"
import "Monica/go-redis/proto"



func main() {
	//获取host
	//获取客户端cmd命令行输入的ip和端口
	var ip = flag.String("ip", "localhost", "redis服务端IP")
	var port = flag.String("port", "6380", "redis服务端PORT")

	flag.Parse()
	host := *ip + ":" + *port
	commandStart := host + "> "

	reader := bufio.NewReader(os.Stdin)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	CheckError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	CheckError(err)

	defer conn.Close()

	
	for {
		fmt.Print(commandStart)
		inputCmd, _ := reader.ReadString('\n')
		inputCmd = strings.Replace(inputCmd, "\n", "", -1)
		send2Server(inputCmd, conn)

		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		CheckError(err)

		if n == 0 {
			fmt.Println(commandStart, "nil")
		}else {
			fmt.Println(string(buff))
		}
	}

}


func send2Server(inputCmd string, conn net.Conn) (n int, err error) {
	encodeCommand, e := proto.EncodeCmd(inputCmd)
	if e != nil {
		return 0, e
	}

	n, err = conn.Write(encodeCommand)
	return n, err
}

func CheckError(err error) {
	if err != nil {
		log.Println("err ", err.Error())
		os.Exit(1)
	}
}