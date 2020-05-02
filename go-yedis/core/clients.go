package core

import (
	"Monica/go-yedis/encrypt"
	"bytes"
	"errors"
)

//
//RedisClient原结构体地址：https://github.com/antirez/redis/blob/30724986659c6845e9e48b601e36aa4f4bca3d30/src/server.h#L765
type YedisClients struct {
	Name string //客户端名称
	Argc int // 当前执行命令的参数的个数
	Argv []*YedisObject //当前执行命令的参数
	Db *YedisDb //指向当前选择数据库的指针
	QueryBuf string //积累客户端查询的缓冲区, 暂用string，后更新用sds
	Reply string //需要发送回客户端的回复信息
	Cmd *YedisCommand //待执行的客户端命令
	LastCommand *YedisCommand //上一个执行的
	Authenticated int //认证状态 0：未认证 1：已认证
	Flags int //客户端状态标志
}


// ProcessInputBuffer 处理客户端请求信息
func (c *YedisClients) ProcessInputBuffer() error {
	decoder := encrypt.NewDecoder(bytes.NewReader([]byte(c.QueryBuf)))
	if resp, err := decoder.DecodeMultiBulk(); err == nil {
		c.Argc = len(resp)
		c.Argv = make([]*YedisObject, c.Argc)
		for k, s := range resp {
			//TODO 写的时候总感觉多此一举，有string导入了还要再创建一个sds，仅仅是模仿Redis的结构，后续考虑按照Go的数据结构来改变
			c.Argv[k] = CreateSdsObject(OBJ_ENCODING_RAW, string(s.Value))
		}
		return nil
	}
	return errors.New("ProcessInputBuffer failed")
}