package sds

import (
	"Monica/go-yedis/core"
	"Monica/go-yedis/utils"
	"strconv"
)

//expire命令   expireat [key] [second] 用于将key对应的值的到期时间设置为倒数second秒
//设置成功返回1，不成功0
func ExpireCommand(c *core.YedisClients, s *core.YedisServer) {

	// 将参数二转时间戳再调用pexpireatCommand
	if second, err := strconv.Atoi(c.Argv[2].Ptr.(core.Sdshdr).Buf); err == nil {
		timestamp := utils.CurrentTimeMillis() + (second * 1000)
		c.Argv[2] = core.CreateSdsObject(core.OBJ_ENCODING_INT, strconv.Itoa(timestamp))
	}
	PexpireatCommand(c, s)

}
