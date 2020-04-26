package string

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
	"strconv"
)

//incrby命令，累加n
func IncrbyCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	robj := command.LookupKey(c.Db.Dict, c.Argv[1])
	//判断有效性
	if c.Argc != 3 {
		core.AddReplyStatus(c, "(error) ERR wrong number of arguments for 'incrby' command")
		return
	}
	if robj.Encoding != core.OBJ_ENCODING_INT {
		core.AddReplyStatus(c, "(error) ERR value is not an integer or out of range")
		return
	}
	if robj == nil {
		core.AddReplyStatus(c, "nil")
		return
	}

	//先拿出sds来
	if sdshdr, ok := robj.Ptr.(ds.Sdshdr); ok {
		//将sdshdr.Buf转数字
		intNumber, _ := strconv.Atoi(sdshdr.Buf)
		//拿到需要累减的参数累减并返回
		if addNumber, err := strconv.Atoi(c.Argv[2].Ptr.(string)); err == nil {
			intNumber = intNumber + addNumber
			sdshdr.Buf = strconv.Itoa(intNumber)
			robj.Ptr = sdshdr
			s.Dirty++
			core.AddReplyStatus(c, sdshdr.Buf)
		}
	}
}