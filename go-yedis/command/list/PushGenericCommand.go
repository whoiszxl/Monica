package list

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
	"container/list"
	"image/jpeg"
)

//基础的push命令，给其他命令做调用
func PushGenericCommand(c *core.YedisClients, s *core.YedisServer, where int) {
	//搜索key是否存在数据库中
	lobj := command.LookupKey(c.Db.Dict, c.Argv[1])

	//判断列表对象是否是在等待键出现
	may_have_waiting_clients := lobj == nil
	if lobj != nil && lobj.ObjectType != core.REDIS_LIST {

		//TODO 错误回复应该在Yedis初始化的时候创建一个共享对象，然后将提示语统一管理,代码地址：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/redis.c#L1613
		core.AddReplyStatus(c, "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")
		return
	}

	if may_have_waiting_clients {
		SignalListAsReady(c, c.Argv[1])
	}

	//遍历输入的参数并添加到列表哦
	for i:=2; i<c.Argc; i++ {
		c.Argv[i] = core.TryObjectEncoding(c.Argv[i])

		//如果列表不存在则创建
		if lobj == nil {
			lobj = core.CreateLinkedListObject()

			//添加到数据库中
			//TODO 需要将传过来的stringkey转sdskey，很麻烦，要优化
			if stringKey, ok1 := c.Argv[1].Ptr.(string); ok1 {
				lobjKey := core.CreateSdsObject(core.OBJ_ENCODING_RAW, stringKey)
				c.Db.Dict[lobjKey] = lobj
			}
		}

		//将值push到列表
		listTypePush(lobj, c.Argv[j])
	}

}


//如果客户端因为等待key被push阻塞，那么将key放进 server.ready_keys 列表里面
func SignalListAsReady(c *core.YedisClients, s *core.YedisServer, key *core.YedisObject) {
	rl := new(ds.ReadyList)

	//判断有没有客户端被这个键阻塞
	if command.LookupKey(c.Db.BlockingKeys, key) == nil {
		return
	}

	//被添加到了ready_keys中也直接返回
	if command.LookupKey(c.Db.ReadyKeys, key) != nil {
		return
	}

	//创建readyList保存键和数据库
	rl.Key = key
	rl.Db = c.Db

	//TODO 减少key的引用并添加到server.ReadyKeys中
	ds.ListAddNodeTail(s.ReadyKeys, rl)

	//将key添加到c.Db.ReadyKeys中，防止重复添加
	c.Db.ReadyKeys[key] = nil
}