package hash

import (
	"Monica/go-yedis/core"
)

func HsetCommand(c *core.YedisClients, s *core.YedisServer) {
	//查找hash对象是否存在
	o := HashTypeLookupWriteOrCreate(c, c.Argv[1])
	if o == nil {
		core.AddReplyError(c, "get or create hash fail")
		return
	}

	update := HashTypeSet(o, c.Argv[2], c.Argv[3])
	if update == 0 {
		core.AddReplyStatus(c, "(integer) 0")
	}else {
		core.AddReplyStatus(c, "(integer) 1")
	}

	//发送信号，发送事件通知 signalModifiedKey notifyKeyspaceEvent

	s.Dirty++

}

//将key-value键值对添加到o的hash中
//key如果存在，则覆盖，不存在则创建
//返回0表示元素已存在，做更新操作，返回1则是添加操作
//Redis提供了ziplist和hashtable来存储，此处简略，只用hashtable来存储
func HashTypeSet(o *core.YedisObject, key *core.YedisObject, value *core.YedisObject) int {

	if o.Encoding == core.OBJ_ENCODING_HT {
		result := core.DictReplace(o.Ptr.(*core.DictHt), key, value)
		if result == 1 {
			//添加操作，增加已用
			o.Ptr.(*core.DictHt).Used++
		}
		return result
	}
	return -1
}



//查找数据库的key的value值，如果不存在则新创建一个
func HashTypeLookupWriteOrCreate(c *core.YedisClients, key *core.YedisObject) *core.YedisObject {

	o := core.LookupKey(c.Db.Dict, c.Argv[1])

	if o == nil {
		//hash不存在
		o = core.CreateHashObject()
		core.DbAdd(c.Db, key, o)
	}else {
		//hash存在,校验类型
		if o.ObjectType != core.REDIS_HASH {
			core.AddReplyError(c, "wrong type err")
			return nil
		}
	}

	return o
}