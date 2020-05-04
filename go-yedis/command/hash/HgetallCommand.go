package hash

import (
	"Monica/go-yedis/core"
)

func HgetallCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找hash对象是否存在
	o := core.LookupKey(c.Db.Dict, c.Argv[1])
	if o == nil {
		core.AddReplyError(c, "(nil)")
		return
	}

	if o.ObjectType != core.REDIS_HASH {
		core.AddReplyError(c, "type err")
		return
	}

	result := ""

	ht := o.Ptr.(*core.DictHt)
	//遍历哈希表，然后全部输出，此处不做批量回复，只添加到一个字符串中
	for _, v := range ht.Table {
		//遍历hash表中的单链表
		iterator := core.DictEntryGetIterator(v)
		for {
			current := core.DictEntryNext(iterator)
			if current == nil {
				break
			}

			//简化回复
			key := current.Key.(core.Sdshdr).Buf
			value := current.Value.(core.Sdshdr).Buf
			kvDict := key + ":" + value
			result = result + "[" + kvDict + "], "
		}
	}

	core.AddReplyStatus(c, "(result) " + result)
}

