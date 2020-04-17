package sds

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"strconv"
)

//pexpireat命令   pexpireat [key] [timestamp] 用于将key对应的值的到期时间为timestamp的毫秒级时间戳
//设置成功返回1，不成功0
func PexpireatCommand(c *core.YedisClients, s *core.YedisServer) {

	//db键空间中有key的话就拿到key对象
	robjKey := command.GetKeyObj(c.Db.Data, c.Argv[1])

	if robjKey != nil {
		//不为空，则在过期db空间里关联键和过期时间
		if timestamp, err := strconv.Atoi(c.Argv[2].Ptr.(string)); err == nil {
			c.Db.Expires[robjKey] = timestamp
			core.AddReplyStatus(c, "1")
		}
	}else {
		core.AddReplyStatus(c, "0")
	}
}
