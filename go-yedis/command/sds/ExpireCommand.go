package sds

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
)

//expire命令   expire [key] [second] 用于将key对应的值的有效期为second秒
func ExpireCommand(c *core.YedisClients, s *core.YedisServer) {

	robj := command.LookupKey(c.Db.Data, c.Argv[1])
	if robj != nil {
		if sdshdr, ok2 := robj.Ptr.(ds.Sdshdr); ok2 {
			core.AddReplyStatus(c, sdshdr.Buf)
		} else {
			core.AddReplyStatus(c, "nil")
		}
	} else {
		core.AddReplyStatus(c, "nil")
	}
}
