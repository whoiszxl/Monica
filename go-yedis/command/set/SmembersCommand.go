package set

import "Monica/go-yedis/core"

//获取集合set中的所有元素列表
func SmembersCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	set := core.LookupKey(c.Db.Dict, c.Argv[1])
	if set != nil {
		htTable := set.Ptr.(*core.DictHt).Table
		result := ""
		for k := range htTable {
			result = result + k.(core.Sdshdr).Buf + ","
		}

		core.AddReplyStatus(c, result)
		return
	}
	core.AddReplyStatus(c, "(integer) 0")
	return
}