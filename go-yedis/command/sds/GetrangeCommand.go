package sds

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
	"strconv"
)

//getrange命令
func GetrangeCommand(c *core.YedisClients, s *core.YedisServer) {
	robj := command.LookupKey(c.Db.Data, c.Argv[1])

	//判断参数有效性
	start, ok := c.Argv[2].Ptr.(string)
	end, ok2 := c.Argv[3].Ptr.(string)
	if !ok || !ok2 {
		core.AddReplyStatus(c, "(error) ERR value is not an integer or out of range")
		return
	}

	//参数转整型
	startNum, err1 := strconv.Atoi(start)
	endNum, err2 := strconv.Atoi(end)
	if err1 != nil || err2 != nil {
		core.AddReplyStatus(c, "(error) ERR value is not an integer or out of range")
		return
	}

	if startNum > endNum {
		core.AddReplyStatus(c, "(error) ERR end must > start")
		return
	}

	if robj != nil {
		if sdshdr, ok := robj.Ptr.(ds.Sdshdr); ok {
			cutStr := sdshdr.Buf[startNum:endNum+1]
			core.AddReplyStatus(c, cutStr)
		}else {
			core.AddReplyStatus(c, "nil")
		}
	}else {
		core.AddReplyStatus(c, "nil")
	}
}