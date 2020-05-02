package list

import (
	"Monica/go-yedis/core"
	"strconv"
)

func LinsertCommand(c *core.YedisClients, s *core.YedisServer) {
	var where int
	//查找list是否存在
	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		//获取插入的方向
		direction := c.Argv[2].Ptr.(core.Sdshdr).Buf
		if direction == "after" {
			where = core.LIST_TAIL
		}else if direction == "before" {
			where = core.LIST_HEAD
		}else {
			core.AddReplyStatus(c, "syntax err")
			return
		}

		//迭代整个list寻找输入的节点并插入输入的节点
		linkedlist := robj.Ptr.(*core.LinkedList)
		inputKey := c.Argv[3].Ptr.(core.Sdshdr).Buf
		findNode := core.ListSearchKey(linkedlist, inputKey)
		inputValue := c.Argv[4]
		core.ListInsertNode(linkedlist, findNode, inputValue, where)


		// signalModifiedKey notifyKeyspaceEvent
		s.Dirty++

		core.AddReplyStatus(c, "(integer) " + strconv.Itoa(linkedlist.Len))
	} else {
		core.AddReplyStatus(c, "(integer) 0")
	}
}