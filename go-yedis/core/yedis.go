package core

import (
	"Monica/go-yedis/encrypt"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
)
//在用户发送命令过来的时候建立客户端连接
func (s *YedisServer) CreateClient() (c *YedisClients) {
	c = new(YedisClients)
	c.Name = string(rand.Intn(10))
	//c.Db = s.yedisDb[0]
	c.Argv = make([]*YedisObject, 5)
	c.Argc = 5
	c.QueryBuf = ""
	c.Reply = ""
	return c
}

//通过connection连接获取客户端请求的命令信息并封装到Client对象中
func (c *YedisClients) ReadCommandFromClient(conn net.Conn) error {
	defer conn.Close()
	buff := make([]byte, 512)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("conn.Read err!=nil", err, "---len---", n, conn)
		conn.Close()
		return err
	}
	c.QueryBuf = string(buff)
	return nil
}

//解密用户客户端发来的加密信息，并将信息存入client中的Argv中，结构为：[0: "set", 1: "name", 2: "www"]
func (c *YedisClients) ProcessCommandInfo() error {
	decoder := encrypt.NewDecoder(bytes.NewReader([]byte(c.QueryBuf)))
	if response, err := decoder.DecodeMultiBulk(); err == nil {
		c.Argc = len(response)
		c.Argv = make([]*YedisObject, c.Argc)
		for count, resp := range response {
			c.Argv[count] = CreateObject(OBJ_STRING, string(resp.Value))
		}
		return nil
	}
	return errors.New("ProcessCommandInfo error")
}

//传入client，执行client中的命令
func (s *YedisServer) ExecuteCommand(c *YedisClients) {

	commandName, ok := c.Argv[0].Ptr.(string)
	if !ok {
		log.Println("error cmd")
		os.Exit(1)
	}

	cmd := LookupCommand(commandName, s)
	if cmd != nil {
		c.Cmd = cmd
		call(c, s)
	}else {
		AddReplyError(c, fmt.Sprintf("(error) ERR unknown command '%s'", commandName))
	}

}

//执行Client中的命令
func call(c *YedisClients, s *YedisServer) {
	dirty := s.Dirty
	c.Cmd.CommandProc(c, s)
	dirty = s.Dirty - dirty

	//判断是否需要aof，开启了则将命令写入server的aofBuff缓冲区
	if s.AofEnabled == ENABLE {
		newStr := append(s.AofBuf, c.QueryBuf)

	}
}


// 查找命令是否支持
func LookupCommand(name string, s *YedisServer) *YedisCommand {
	if cmd, ok := s.Commands[name]; ok {
		return cmd
	}
	return nil
}