package db

import (
	"Monica/go-yedis/core"
)

//获取所有的key
func KeysCommand(c *core.YedisClients, s *core.YedisServer) {
	pattern := c.Argv[1].Ptr.(string)

	result := ""
	if pattern == "*" {
		//暂只实现查询所有的功能
		for k, _ := range c.Db.Dict {
			keyStr := k.Ptr.(core.Sdshdr).Buf
			result = result + keyStr + ","
		}

	}
	if result != "" {
		core.AddReplyStatus(c, result)
	}else {
		core.AddReplyStatus(c, "empty database")
	}

}
