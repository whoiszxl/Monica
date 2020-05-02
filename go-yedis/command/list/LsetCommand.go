package list

import (
	"Monica/go-yedis/core"
	"strconv"
)

func LsetCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找list是否存在
	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		//获取用户输入的下标
		index := c.Argv[2].Ptr.(core.Sdshdr).Buf
		indexNum, _ := strconv.Atoi(index)
		//获取用户需要设置的值
		value := c.Argv[3]
		//将从数据库查询的记录强转为双向链表,并对参数做校验
		linkedList := robj.Ptr.(*core.LinkedList)
		if indexNum > linkedList.Len || indexNum < 0 {
			core.AddReplyError(c, "(error) ERR index out of range")
			return
		}
		//查询出index节点并设置新的值
		listIndexNode := core.ListIndex(linkedList, indexNum)
		listIndexNode.Value = value
		core.AddReplyStatus(c, "OK")
	} else {
		core.AddReplyStatus(c, "(error) ERR no such key")
	}
}