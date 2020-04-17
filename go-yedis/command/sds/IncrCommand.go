package sds

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"strconv"
)

//incr命令，累加1
func IncrCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	robj := command.LookupKey(c.Db.Data, c.Argv[1])
	//判断类型有效性
	if robj == nil {
		core.AddReplyStatus(c, "nil")
	}
	if robj.Encoding != core.OBJ_ENCODING_INT {
		core.AddReplyStatus(c, "(error) ERR value is not an integer or out of range")
		return
	}

	if intFloat, ok := robj.Ptr.(float64); ok {
		//在ptr中可以直接拿到int類型
		intFloat = intFloat + 1
		robj.Ptr = intFloat
		core.AddReplyStatus(c, strconv.FormatFloat(intFloat, 'f', -1, 64))
	}
}
