package list

import (
	"Monica/go-yedis/core"
	"strconv"
)

func LindexCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找list是否存在
	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		index, err := strconv.Atoi(c.Argv[2].Ptr.(string))
		if err != nil {
			core.AddReplyError(c, "(nil)")
			return
		}

		listNode := core.ListIndex(robj.Ptr.(*core.LinkedList), index)
		if index >= robj.Ptr.(*core.LinkedList).Len {
			core.AddReplyError(c, "(nil)")
			return
		}
		core.AddReplyStatus(c, listNode.Value.(string))
	} else {
		core.AddReplyStatus(c, "(integer) 0")
	}
}