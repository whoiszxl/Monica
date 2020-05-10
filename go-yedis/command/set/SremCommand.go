package set

import (
	"Monica/go-yedis/core"
	"strconv"
)

//从集合中删除一个元素
func SremCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	set := core.LookupKey(c.Db.Dict, c.Argv[1])
	if set == nil {
		core.AddReplyStatus(c, "set is nil")
	}else {
		if set.ObjectType != core.REDIS_SET {
			core.AddReplyStatus(c, "wrong type err")
			return
		}
	}

	deleted := 0
	for j:=2; j<c.Argc; j++ {
		if setTypeDelete(set, c.Argv[j].Ptr.(core.Sdshdr)) {
			deleted++
		}
	}

	if deleted > 0 {
		//signalModifiedKey(c,c->db,c->argv[1]);
		//notifyKeyspaceEvent(NOTIFY_SET,"sadd",c->argv[1],c->db->id);

		//增加HT的used引用,直接用golang原生len方法，也是O(1)复杂度
		set.Ptr.(*core.DictHt).Used = len(set.Ptr.(*core.DictHt).Table)
	}

	s.Dirty += deleted
	core.AddReplyStatus(c, "(integer) " + strconv.Itoa(deleted))
}


func setTypeDelete(subject *core.YedisObject, value core.Sdshdr) bool {
	delete(subject.Ptr.(*core.DictHt).Table, value)
	return true
}