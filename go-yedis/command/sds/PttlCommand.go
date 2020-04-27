package sds

import (
	"Monica/go-yedis/core"
	"Monica/go-yedis/utils"
	"strconv"
)

//pttl命令   pttl [key] 用于获取key的失效时间，单位毫秒
func PttlCommand(c *core.YedisClients, s *core.YedisServer) {

	validTimestamp := GetPttlTime(c, s)
	core.AddReplyStatus(c, strconv.Itoa(validTimestamp))
}


//获取毫秒级别的失效时间
func GetPttlTime(c *core.YedisClients, s *core.YedisServer) int {
	//db键空间中有key的话就拿到key对象
	robjKey := core.GetKeyObj(c.Db.Dict, c.Argv[1])
	if robjKey == nil {
		return -2
	}

	//从expireDict中拿到失效时间戳
	expireTimestamp := c.Db.Expires[robjKey]
	if expireTimestamp == 0 {
		return -1
	}

	//键值对的有效时间 = 过期日期时间戳 - 当前时间戳
	currentMillis := utils.CurrentTimeMillis()
	validTimestamp := expireTimestamp - currentMillis

	if validTimestamp < 0 {
		return -2
	}
	return validTimestamp
}
