package hash

import (
	"Monica/go-yedis/core"
	"strconv"
)

func HlenCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找hash对象是否存在
	o := core.LookupKey(c.Db.Dict, c.Argv[1])
	if o == nil {
		core.AddReplyError(c, "get or create hash fail")
		return
	}

	if o.ObjectType != core.REDIS_HASH {
		core.AddReplyError(c, "type err")
		return
	}

	ht := o.Ptr.(*core.DictHt)
	core.AddReplyStatus(c, "(integer) " + strconv.Itoa(ht.Used))
}
