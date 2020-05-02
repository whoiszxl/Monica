package list

import (
	"Monica/go-yedis/core"
)

//在LinkedList中弹出一个元素并返回
func PopGenericCommand(c *core.YedisClients, s *core.YedisServer, where int) {
	//搜索key是否存在数据库中，以及判断value是否是list
	o := core.LookupKey(c.Db.Dict, c.Argv[1])
	if o == nil || o.ObjectType != core.REDIS_LIST {
		core.AddReplyStatus(c, "(nil)")
		return
	}

	//获取头部或尾部节点并删除其
	value := listTypePop(c, o, where)
	if value == nil {
		core.AddReplyError(c, "nil")
		return
	}

	//判断现在的value是否还有元素，没有则删除键
	if o.Ptr.(*core.LinkedList).Len == 0 {
		//TODO 遍历整个Dict拿到key对象，然后再通过对象删除，感觉还有很大优化空间，O(1)变O(N)了，很不给力
		keyObj := core.GetKeyObj(c.Db.Dict, c.Argv[1])
		delete(c.Db.Dict, keyObj)
	}

	core.AddReplyStatus(c, value.Value.Ptr.(core.Sdshdr).Buf)

	s.Dirty++
}



func LpopCommand(c *core.YedisClients, s *core.YedisServer) {
	PopGenericCommand(c, s, core.LIST_HEAD)
}


func RpopCommand(c *core.YedisClients, s *core.YedisServer) {
	PopGenericCommand(c, s, core.LIST_TAIL)
}


func listTypePop(c *core.YedisClients, subject *core.YedisObject, where int) *core.ListNode {
	var node *core.ListNode
	linkedlist := subject.Ptr.(*core.LinkedList)
	if subject.Encoding == core.OBJ_ENCODING_LINKEDLIST {
		if where == core.LIST_HEAD {
			//获取头部的节点并移除头部节点
			node = linkedlist.Head
			core.ListDelNode(linkedlist, node)
		}else {
			node = linkedlist.Tail
			core.ListDelNode(linkedlist, node)
		}
		return node
	}
	return nil
}
