package string

import (
	"Monica/go-yedis/core"
	"Monica/go-yedis/utils"
	"strconv"
)

//pexpire命令   pexpireat [key] [millis] 用于将key对应的值的到期时间设置为倒数millis毫秒
//设置成功返回1，不成功0
func PexpireCommand(c *core.YedisClients, s *core.YedisServer) {

	// 将参数二转时间戳再调用pexpireatCommand
	if millis, err := strconv.Atoi(c.Argv[2].Ptr.(string)); err == nil {
		timestamp := utils.CurrentTimeMillis() + millis
		c.Argv[2].Ptr = strconv.Itoa(timestamp)
	}
	PexpireatCommand(c, s)
}
