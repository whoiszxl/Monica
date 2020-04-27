package core

import (
	"Monica/go-yedis/encrypt"
	"Monica/go-yedis/utils"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
)

const (
	REDIS_CALL_NONE      = 0
	REDIS_CALL_SLOWLOG   = 1
	REDIS_CALL_STATS     = 2
	REDIS_CALL_PROPAGATE = 4
	REDIS_CALL_FULL      = (REDIS_CALL_SLOWLOG | REDIS_CALL_STATS | REDIS_CALL_PROPAGATE)

	REDIS_AOF_OFF = 0
	REDIS_AOF_ON  = 1

	REDIS_PROPAGATE_NONE = 0
	REDIS_PROPAGATE_AOF  = 1
	REDIS_PROPAGATE_REPL = 2
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
			//判断客户端传来的Value是什么类型 (int sds) ....不判断了，string放进去就完事了
			c.Argv[count] = CreateObject(REDIS_STRING, OBJ_ENCODING_RAW, string(resp.Value))
		}
		return nil
	}
	return errors.New("ProcessCommandInfo error")
}

//传入client，执行client中的命令
func (s *YedisServer) ExecuteCommand(c *YedisClients) {

	//quit命令单独处理
	if c.Argv[0].Ptr.(string) == "quit" {
		AddReplyStatus(c, "bye bye")
		os.Exit(1)
	}

	//校验传入的命令有效性
	commandName, ok := c.Argv[0].Ptr.(string)
	if !ok {
		log.Println("error cmd")
		return
	}

	//查找Yedis是否支持此命令
	cmd := LookupCommand(commandName, s)

	//校验参数个数是否正确
	if (cmd.Arity > 0 && cmd.Arity != c.Argc) || (c.Argc < -cmd.Arity) {
		AddReplyError(c, fmt.Sprintf("(error) wrong number of arguments for '%s' command", cmd.Name))
		return
	}

	//密码校验
	if s.Requirepass != "" && c.Authenticated != 1 {
		AddReplyError(c, "NO AUTH")
		return
	}

	//TODO 集群处理
	//TODO 如果设置了最大内存，检查是否超过限制，超过了则去删除过期键来释放内存

	if cmd != nil {
		c.Cmd = cmd
		call(c, s, REDIS_CALL_FULL)
	}else {
		AddReplyError(c, fmt.Sprintf("(error) ERR unknown command '%s'", commandName))
		return
	}
}

//执行Client中的命令
func call(c *YedisClients, s *YedisServer, flags int) {
	//保留旧Dirty值，并在执行命令后计算差值
	dirty := s.Dirty
	startTime := utils.CurrentTimeMicrosecond()
	c.Cmd.CommandProc(c, s)
	endTime := utils.CurrentTimeMicrosecond()
	//计算命令消耗的时间和Dirty差值,微妙级别
	duration := endTime - startTime
	dirty = s.Dirty - dirty

	//TODO LUA相关处理

	//是否需要将命令放到慢日志（Slowlog）中
	if flags & REDIS_CALL_SLOWLOG != 0{
		SlowlogPushEntryIfNeeded(c.Argv, c.Argc, duration)
	}

	//是否需要更新统计信息，执行命令耗费的微妙数和调用次数
	if flags & REDIS_CALL_STATS != 0 {
		c.Cmd.Microseconds += int64(duration)
		c.Cmd.Calls++
	}

	//里面判断挺麻烦的，直接propagate传播一下，简化一下
	if flags & REDIS_CALL_PROPAGATE != 0 {
		if dirty > 0 {
			propagate(s, c, flags)
		}
	}

}

//将命令传播到aof， slave的话先省略了
func propagate(server *YedisServer, client *YedisClients, flags int) {
	if server.AofState != REDIS_AOF_OFF && flags & REDIS_PROPAGATE_AOF > 0 {
		//传播到AOF
		FeedAppendOnlyFile(server, client)
	}

	if flags & REDIS_PROPAGATE_REPL > 0 {
		//TODO 传播到slave暂不实现
		//replicationFeedSlaves(server.slaves,dbid,argv,argc);
	}
}


// 查找命令是否支持
func LookupCommand(name string, s *YedisServer) *YedisCommand {
	if cmd, ok := s.Commands[name]; ok {
		return cmd
	}
	return nil
}