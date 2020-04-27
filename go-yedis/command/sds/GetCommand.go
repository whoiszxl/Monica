package sds

import (
	"Monica/go-yedis/core"
)

//get命令
func GetCommand(c *core.YedisClients, s *core.YedisServer) {

	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		if sdshdr, ok2 := robj.Ptr.(core.Sdshdr); ok2 {
			core.AddReplyStatus(c, sdshdr.Buf)
		} else {
			core.AddReplyStatus(c, "nil")
		}
	} else {
		core.AddReplyStatus(c, "nil")
	}
}
