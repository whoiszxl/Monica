package sds

import (
	"Monica/go-yedis/core"
	"strconv"
)

//expireat命令   expireat [key] [timestamp] 用于将key对应的值的到期时间为timestamp的秒级时间戳
//设置成功返回1，不成功0
func ExpireatCommand(c *core.YedisClients, s *core.YedisServer) {

	//db键空间中有key的话就拿到key对象
	robjKey := core.GetKeyObj(c.Db.Dict, c.Argv[1])

	if robjKey != nil {
		//不为空，则在过期db空间里关联键和过期时间
		if timestamp, err := strconv.Atoi(c.Argv[2].Ptr.(core.Sdshdr).Buf); err == nil {
			//获取的是秒级时间戳，需要乘1000再存入
			c.Db.Expires[robjKey] = timestamp * 1000
			core.AddReplyStatus(c, "1")
		}
	}else {
		core.AddReplyStatus(c, "0")
	}
}