package db

import (
	"Monica/go-yedis/core"
	"strconv"
)

//如果 key 已经存在并且是一个字符串， APPEND 命令将指定的 value 追加到该 key 原来值（value）的末尾
func SelectCommand(c *core.YedisClients, s *core.YedisServer) {
	id, _ := strconv.Atoi(c.Argv[1].Ptr.(string))
	if id < 0 || id >= s.DbNum {
		core.AddReplyStatus(c, "invalid DB index")
		return
	}

	if core.SelectDb(c, s, id) == core.REDIS_OK {
		core.AddReplyStatus(c, "OK")
	}else {
		core.AddReplyStatus(c, "invalid DB index")
	}

}
