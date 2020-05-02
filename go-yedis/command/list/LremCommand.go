package list

import (
	"Monica/go-yedis/core"
	"strconv"
)

func LremCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找list是否存在
	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		count, err := strconv.Atoi(c.Argv[2].Ptr.(core.Sdshdr).Buf)
		if err != nil {
			core.AddReplyError(c, "(nil)")
			return
		}

		//TODO 暂时只做0的处理
		//var where int
		//if count >= 0 {
		//	where = core.LIST_HEAD
		//}else {
		//	where = core.LIST_TAIL
		//}

		remValue := c.Argv[3].Ptr.(core.Sdshdr).Buf
		counter := 0
		list := robj.Ptr.(*core.LinkedList)
		if count == 0 {
			iter := core.ListGetIterator(list, core.AL_START_HEAD)

			for {
				node := core.ListNext(iter)
				if node == nil {
					break
				}
				if remValue == node.Value.Ptr.(core.Sdshdr).Buf {
					core.ListDelNode(list, node)
					counter++
				}
			}
		}

		core.AddReplyStatus(c, "(integer) " + strconv.Itoa(counter))
	} else {
		core.AddReplyStatus(c, "(integer) 0")
	}
}