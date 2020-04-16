package command

import (
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
)

// set命令
// Redis源码中key用的是robj,这里精简一下直接采用go原生string
func SetCommand(c *core.YedisClients, s *core.YedisServer) {
	//校验有效性
	if !checkParam(c) {
		return
	}

	//获取键值对
	robjKey := c.Argv[1]
	robjValue := c.Argv[2]

	//判断是否是字符串，是则设置到Db的Data中
	if stringKey, ok1 := robjKey.Ptr.(string); ok1 {
		if stringValue, ok2 := robjValue.Ptr.(string); ok2 {
			//创建一个sdshdr保存到字典中
			robjSds := ds.Sdshdr{Len:uint(len(stringValue)), Free:0, Buf:stringValue}
			c.Db.Data[stringKey] = core.CreateObject(core.OBJ_STRING, robjSds)
		}
	}

	s.Dirty++
	core.AddReplyStatus(c, "OK")
}


//校验参数是否有效
func checkParam(c *core.YedisClients) bool{
	if c.Argc < 3 {
		core.AddReplyError(c, "(error) ERR wrong number of arguments for 'set' command")
		return false
	}
	if c.Argc > 3 {
		core.AddReplyError(c, "(error) ERR syntax error")
		return false
	}
	return true
}