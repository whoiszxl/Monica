package sds

import (
	"Monica/go-yedis/core"
	"strconv"
	"time"
)

//expire命令   expireat [key] [second] 用于将key对应的值的到期时间设置为倒数second秒
//设置成功返回1，不成功0
func ExpireCommand(c *core.YedisClients, s *core.YedisServer) {

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
