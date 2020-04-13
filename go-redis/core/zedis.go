package core

import (
	"fmt"
	"log"
	"net"
	"os"
	"bytes"
	"errors"
	"Monica/go-redis/proto"
	"Monica/go-redis/persistence"
)

// 服务端连接之后创建一个和用户端对应的服务Client
type Client struct {
	Cmd *ZedisCommand
	Argv []*ZedisObject
	Argc int
	Db *ZedisDb
	QueryBuf string
	Buf string
}

// Server 服务端实例结构体
type Server struct {
	Db               []*ZedisDb
	DbNum            int
	Start            int64
	Port             int32
	RdbFilename      string
	AofFilename      string
	NextClientID     int32
	SystemMemorySize int32
	Clients          int32
	Pid              int
	Commands         map[string]*ZedisCommand
	Dirty            int64
	AofBuf           []string
}

type ZedisCommand struct {
	Name string
	Proc cmdFunc
}

type cmdFunc func(c *Client, s *Server)

//get命令，通过参数1去寻找
func GetCommand(c *Client, s *Server) {
	//aa := c.Db.Dict["name"].Ptr.(string)
	o := findKey(c.Db.Dict, c.Argv[1])
	if o != nil {
		addReplyStatus(c, o.Ptr.(string))
	}else {
		addReplyStatus(c, "nil")
	}
}

func findKey(dicts dict, key *ZedisObject) (ret *ZedisObject) {

	//TODO 无法dicts["key"]访问，不知道什么鬼
	for key, val := range dicts {
		if key == "name" {
			return val
		}
	}
	return nil
}

//set命令
func SetCommand(c *Client, s *Server) {
	//set命令参数有效性校验
	if c.Argc != 3 {
		addReply(c, CreateObject(ObjectTypeString, "(error) ERR syntax error"))
		return
	}
	
	//获取键值对
	objKey := c.Argv[1]
	objValue := c.Argv[2]

	//判断是否是字符串，是的话设置到Db的Dict中
	if stringKey, ok1 := objKey.Ptr.(string); ok1 {
		if stringValue, ok2 := objValue.Ptr.(string); ok2 {
			c.Db.Dict[stringKey] = CreateObject(ObjectTypeString, stringValue)
		}
	}
	
	for k, v := range c.Db.Dict {
		fmt.Println(k)
		fmt.Println(v.Ptr.(string))
	}

	s.Dirty++
	addReplyStatus(c, "OK")
}


//使用Go原生数据结构map作为redis中dict结构体
type dict map[string]*ZedisObject

type ZedisDb struct {
	Dict dict
	Expires dict
	ID int32
}



//在用户发送命令过来的时候建立客户端连接
func (s *Server) CreateClient() (c *Client) {
	c = new(Client)
	c.Db = s.Db[0]
	c.Argv = make([]*ZedisObject, 5)
	c.QueryBuf = ""
	return c
}


//通过connection连接获取客户端请求信息并封装到Client对象中
func (c *Client) ReadQueryFromClient(conn net.Conn) (err error) {
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

// 将命令切割放入c.Argv
func (c *Client) ProcessInputBuffer() error {
	decoder := proto.NewDecoder(bytes.NewReader([]byte(c.QueryBuf)))
	if resp, err := decoder.DecodeMultiBulk(); err == nil {
		c.Argc = len(resp)
		c.Argv = make([]*ZedisObject, c.Argc)
		for k, s := range resp {
			c.Argv[k] = CreateObject(ObjectTypeString, string(s.Value))
		}
		return nil
	}
	return errors.New("ProcessInputBuffer failed")
}


//执行命令
func (s *Server) ProcessCommand(c *Client) {
	v := c.Argv[0].Ptr
	name, ok := v.(string)
	if !ok {
		log.Println("error cmd")
		os.Exit(1)
	}

	cmd := findCommand(name, s)
	if cmd != nil {
		c.Cmd = cmd
		call(c, s)
	}else {
		addReplyError(c, fmt.Sprintf("(error) ERR unknown command '%s'", name))
	}
}

func call(c *Client, s *Server) {
	dirty := s.Dirty
	c.Cmd.Proc(c, s)
	dirty = s.Dirty - dirty
	if dirty > 0 {
		persistence.AppendToFile(s.AofFilename, c.QueryBuf)
	}
}

// 查找命令是否支持
func findCommand(name string, s *Server) *ZedisCommand {
	if cmd, ok := s.Commands[name]; ok {
		return cmd
	}
	return nil
}


//添加回复
func addReply(c *Client, o *ZedisObject) {
	c.Buf = o.Ptr.(string)
}

func addReplyStatus(c *Client, s string) {
	r := proto.NewString([]byte(s))
	addReplyString(c, r)
}
func addReplyError(c *Client, s string) {
	r := proto.NewError([]byte(s))
	addReplyString(c, r)
}
func addReplyString(c *Client, r *proto.Resp) {
	if ret, err := proto.EncodeToBytes(r); err == nil {
		c.Buf = string(ret)
	}
}