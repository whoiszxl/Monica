package sds

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
)

//get命令
func GetCommand(c *core.YedisClients, s *core.YedisServer) {

	robj := command.LookupKey(c.Db.Data, c.Argv[1])
	if robj != nil {
		if sdshdr, ok := robj.Ptr.(ds.Sdshdr); ok {
			core.AddReplyStatus(c, sdshdr.Buf)
		} else {
			core.AddReplyStatus(c, "nil")
		}
	} else {
		core.AddReplyStatus(c, "nil")
	}
}
