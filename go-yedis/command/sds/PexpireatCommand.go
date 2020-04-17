package sds

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"strconv"
)

//pexpireat命令   pexpireat [key] [timestamp] 用于将key对应的值的到期时间为timestamp的毫秒级时间戳
//设置成功返回1，不成功0
func PexpireatCommand(c *core.YedisClients, s *core.YedisServer) {

	//判断key是否在键空间db中
	robj := command.LookupKey(c.Db.Data, c.Argv[1])
	if robj != nil {
		//不为空，则在过期db空间里关联键和过期时间
		if timestamp, err := strconv.Atoi(c.Argv[2].Ptr.(string)); err == nil {
			key := c.Argv[1].Ptr.(*core.YedisObject)
			c.Db.Expires[key] = timestamp
			core.AddReplyStatus(c, "1")
		}
	}else {
		core.AddReplyStatus(c, "0")
	}
}
