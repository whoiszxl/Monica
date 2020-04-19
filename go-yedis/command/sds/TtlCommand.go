package sds

import (
	"Monica/go-yedis/core"
	"strconv"
	"time"
)

//ttl命令   ttl [key] 用于获取key的过期秒数
//设置成功返回1，不成功0
func TtlCommand(c *core.YedisClients, s *core.YedisServer) {

	//直接

	//获取传入的秒数
	if second, err := strconv.Atoi(c.Argv[2].Ptr.(string)); err == nil {
		//将当前的时间加上秒数并传递回原有秒数参数内
		newTime := int(time.Now().UnixNano() / 1000000) + second
		c.Argv[2].Ptr = newTime
		ExpireatCommand(c, s)
	}else {
		core.AddReplyStatus(c, "0")
	}
}
