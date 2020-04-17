package sds

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
	"strconv"
)

//decr命令，累减1
func DecrCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	robj := command.LookupKey(c.Db.Data, c.Argv[1])
	//判断有效性
	if c.Argc != 2 {
		core.AddReplyStatus(c, "(error) ERR wrong number of arguments for 'decr' command")
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
		intNumber = intNumber - 1
		sdshdr.Buf = strconv.Itoa(intNumber)
		robj.Ptr = sdshdr
		core.AddReplyStatus(c, sdshdr.Buf)
	}
}