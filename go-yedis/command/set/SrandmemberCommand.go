package set

import "Monica/go-yedis/core"

//随机获取一个元素
func SrandmemberCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	set := core.LookupKey(c.Db.Dict, c.Argv[1])
	if set != nil {
		htTable := set.Ptr.(*core.DictHt).Table

		//从htTable中随机获取一个元素
		key := RandDictMapKey(htTable, false)
		core.AddReplyStatus(c, key)
		return
	}
	core.AddReplyStatus(c, "nil")
	return
}