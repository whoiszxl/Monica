package list

import (
	"Monica/go-yedis/core"
	"strconv"
)

func LrangeCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找list是否存在
	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		//将数据库中查询到的对象强转为链表
		list := robj.Ptr.(*core.LinkedList)
		llen := list.Len

		//设定start和end
		start, err1 := strconv.Atoi(c.Argv[2].Ptr.(string))
		end, err2 := strconv.Atoi(c.Argv[3].Ptr.(string))
		if err1 != nil || err2 != nil {
			core.AddReplyError(c, "(error) ERR value is not an integer or out of range")
			return
		}

		if start < 0 {
			start = llen + start
		}
		if end < 0 {
			end = llen + end
		}
		if start < 0 {
			start = 0
		}
		if start > end || start >= llen {
			core.AddReplyError(c, "(empty list or set)")
			return
		}
		if end >= llen {
			end = llen - 1
		}
		rangelen := (end - start) + 1

		//开始批量回复，TODO 这里暂时不批量回复，先拼接成字符串回复先
		iter := core.ListGetIterator(list, core.AL_START_HEAD)

		result := ""
		for rangelen > 0 {
			node := core.ListNext(iter)
			if node == nil {
				break
			}
			result = result + node.Value.(string) + ","
			rangelen--
		}

		core.AddReplyStatus(c, "(linkedlist) " + result)
	} else {
		core.AddReplyStatus(c, "(empty list or set)")
	}
}