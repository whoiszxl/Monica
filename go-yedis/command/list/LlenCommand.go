package list

import (
	"Monica/go-yedis/core"
	"strconv"
)

func LlenCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找list是否存在
	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		list := robj.Ptr.(*core.LinkedList)
		core.AddReplyStatus(c, "(integer) " + strconv.Itoa(list.Len))
	} else {
		core.AddReplyStatus(c, "(integer) 0")
	}
}