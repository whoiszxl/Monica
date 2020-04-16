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
	if robj != nil {
		if sdshdr, ok := robj.Ptr.(ds.Sdshdr); ok {
			//获取到字符串sdshdr对象,判断它是否能转int
			//TODO 直接覆盖而不是修改内存值，暂有优化空间
			if number, err := strconv.Atoi(sdshdr.Buf); err == nil {
				number = number - 1
				sdshdr.Buf = strconv.Itoa(number)
				robj.Ptr = sdshdr
			}

			core.AddReplyStatus(c, sdshdr.Buf)
		}else {
			core.AddReplyStatus(c, "nil")
		}
	}else {
		core.AddReplyStatus(c, "nil")
	}
}