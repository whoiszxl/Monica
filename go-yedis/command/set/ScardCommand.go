package set

import (
	"Monica/go-yedis/core"
	"strconv"
)

//获取set元素的数量
func ScardCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	set := core.LookupKey(c.Db.Dict, c.Argv[1])
	if set != nil {
		used := set.Ptr.(*core.DictHt).Used
		core.AddReplyStatus(c, "(integer) " + strconv.Itoa(used))
		return
	}
	core.AddReplyStatus(c, "nil")
	return
}