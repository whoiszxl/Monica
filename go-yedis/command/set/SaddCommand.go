package set

import (
	"Monica/go-yedis/core"
	"strconv"
)

//调用dictAdd，新元素为key，nil为值。添加到dict键值对中
//Redis源码：https://github.com/antirez/redis/blob/3.0/src/t_set.c#L250
func SaddCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	set := core.LookupKey(c.Db.Dict, c.Argv[1])
	if set == nil {
		set = setTypeCreate(c.Argv[2].Ptr.(core.Sdshdr))
		core.DbAdd(c.Db, c.Argv[1], set)
	}else {
		if set.ObjectType != core.REDIS_SET {
			core.AddReplyStatus(c, "wrong type err")
			return
		}
	}

	added := 0
	for j:=2; j<c.Argc; j++ {
		if setTypeAdd(set, c.Argv[j].Ptr.(core.Sdshdr)) {
			added++
		}
	}

	if added > 0 {
		//signalModifiedKey(c,c->db,c->argv[1]);
		//notifyKeyspaceEvent(NOTIFY_SET,"sadd",c->argv[1],c->db->id);

		//增加HT的used引用,直接用golang原生len方法，也是O(1)复杂度
		set.Ptr.(*core.DictHt).Used = len(set.Ptr.(*core.DictHt).Table)
	}

	s.Dirty += added
	core.AddReplyStatus(c, "(integer) " + strconv.Itoa(added))

}

//此处需要判断是否创建intSet还是dictSet,此处简化，直接创建dictSet
func setTypeCreate(sdshdr core.Sdshdr) *core.YedisObject{
	return core.CreateSetObject()
}

func setTypeAdd(subject *core.YedisObject, value core.Sdshdr) bool {
	if subject.ObjectType == core.REDIS_SET && subject.Encoding == core.OBJ_ENCODING_HT {
		ht := subject.Ptr.(*core.DictHt)
		ht.Table[value] = nil
		return true
	}
	return false
}