package hash

import (
	"Monica/go-yedis/core"
	"Monica/go-yedis/utils"
)

func HdelCommand(c *core.YedisClients, s *core.YedisServer) {
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

	//获取hash表map
	ht := o.Ptr.(*core.DictHt)
	htTable := ht.Table

	field := c.Argv[2].Ptr.(core.Sdshdr).Buf
	encodingHash := utils.Times33Encoding(field)
	index := encodingHash % core.DEFAULT_HASH_LEN
	entry := htTable[int(index)]

	//如果下标的链表只有一个元素，直接删除
	if entry.Next == nil {
		delete(htTable, int(index))
		core.AddReplyStatus(c, "(integer) 1")
		return
	}

	iterator := core.DictEntryGetIterator(entry)
	var prev *core.DictEntry

	for {
		current := core.DictEntryNext(iterator)
		if current == nil {
			break
		}

		if field == current.Key.(core.Sdshdr).Buf {
			//如果当前prev是nil，则是表头，需要删除表头,直接将map数组下标的元素替换为当前的next
			if prev == nil {
				htTable[int(index)] = current.Next
			}else {
				//删除这个节点，将上一个的节点指向下下个
				prev.Next = current.Next
			}
			core.AddReplyStatus(c, "(integer) 1")
			return
		}
		prev = current
	}

	s.Dirty++
	core.AddReplyStatus(c, "(integer) 0")
}
