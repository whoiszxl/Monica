package hash

import (
	"Monica/go-yedis/core"
	"Monica/go-yedis/utils"
)

func HexistsCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找hash对象是否存在
	o := core.LookupKey(c.Db.Dict, c.Argv[1])
	if o == nil {
		core.AddReplyStatus(c, "(nil)")
		return
	}
	ht := o.Ptr.(*core.DictHt)
	htTable := ht.Table
	field := c.Argv[2].Ptr.(core.Sdshdr).Buf
	encodingHash := utils.Times33Encoding(field)
	index := encodingHash % core.DEFAULT_HASH_LEN
	entry := htTable[int(index)]

	iterator := core.DictEntryGetIterator(entry)

	for {
		current := core.DictEntryNext(iterator)
		if current == nil {
			break
		}

		if field == current.Key.(core.Sdshdr).Buf {
			core.AddReplyStatus(c, "(integer) 1")
			return
		}
	}

	core.AddReplyStatus(c, "(integer) 0")
	return

}