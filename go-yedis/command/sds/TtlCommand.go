package sds

import (
	"Monica/go-yedis/core"
	"strconv"
)

//ttl命令   ttl [key] 用于获取key的失效时间，单位秒
func TtlCommand(c *core.YedisClients, s *core.YedisServer) {

	validTimestamp := GetPttlTime(c, s)
	core.AddReplyStatus(c, strconv.Itoa(validTimestamp / 1000))
}

