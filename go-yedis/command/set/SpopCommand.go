package set

import (
	"Monica/go-yedis/core"
	"math/rand"
)

//随机弹出一个元素
func SpopCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	set := core.LookupKey(c.Db.Dict, c.Argv[1])
	if set != nil {
		htTable := set.Ptr.(*core.DictHt).Table

		//从htTable中随机获取一个元素
		key := RandDictMapKey(htTable, true)
		s.Dirty--
		//TODO 此命令重写，AOF需要将命令转换为srem
		core.AddReplyStatus(c, key)
		return
	}
	core.AddReplyStatus(c, "nil")
	return
}

//从DictMap中随机获取一个key
func RandDictMapKey(m core.DictMap, isDelete bool) string {

	randomNumber := rand.Intn(len(m))

	result := ""
	for k := range m {
		if randomNumber == 0 {
			result = k.(core.Sdshdr).Buf
			if isDelete {
				delete(m, k)
			}
			break
		}
		randomNumber--
	}

	return result
}