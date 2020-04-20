package core

import (
	"Monica/go-yedis/encrypt"
	"Monica/go-yedis/persistence"
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
	c.Db = s.ServerDb[0]
	c.Argv = make([]*YedisObject, 5)
	c.Argc = 5
	c.QueryBuf = ""
	c.Reply = ""
	return c
}

//Redis的定时任务器，每秒钟调用config.hz次，默认是每秒十次
//其中Yedis实现的异步操作如下:
//1. 去Db.ExpireDict下清除过期的key
//2. 更新一些统计信息，比如说内存使用情况，内存最高占用，command的平均调用时间
//3. 调用bgsave备份数据到rdb，或者进行aof重写
//4. 清除超时客户端连接
func ServerCron(loop *AeEventLoop, server *YedisServer) int {

	//数据库相关定时任务
	databasesCron()

	//客户端相关定时任务
	clientCron()

	//记录服务器的内存峰值
	RecordPeakMemory()


	//如果没有后台重写aof和rdb bgsave, 检查是否需要执行bgsave
	//检查更新的数量是否大于配置数量，还有时间是否超过了配置时间
	if server.Dirty >= server.SaveNumber && server.Unixtime - server.LastSaveTime > server.SaveNumber {
		persistence.RdbSaveBackground(server.RdbFileName)
	}

	//TODO 判断是否要执行aof重写
	if server.RdbChildPid == -1 && server.AofChildPid == -1 {
		persistence.RewriteAppendOnlyFileBackground()
	}

	return 1
}

//数据库相关处理的执行函数
func databasesCron() {
	log.Println("开始执行databasesCron")
}

//客户端相关定时任务
func clientCron() {

}

func RecordPeakMemory() {

}

//通过connection连接获取客户端请求的命令信息并封装到Client对象中
func (c *YedisClients) ReadCommandFromClient(conn net.Conn) error {
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
			//判断客户端传来的Value是什么类型 (int string) ....不判断了，string放进去就完事了
			c.Argv[count] = CreateObject(OBJ_STRING, OBJ_ENCODING_RAW, string(resp.Value))
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
		s.AofBuf = append(s.AofBuf, c.QueryBuf)
	}
}


// 查找命令是否支持
func LookupCommand(name string, s *YedisServer) *YedisCommand {
	if cmd, ok := s.Commands[name]; ok {
		return cmd
	}
	return nil
}