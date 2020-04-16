package sds

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
	"strings"
)

//mget命令
func MgetCommand(c *core.YedisClients, s *core.YedisServer) {


	//循环获取
	var result = make([]string, c.Argc - 1)
	for i := 1; i < c.Argc; i++ {
		robj := command.LookupKey(c.Db.Data, c.Argv[i])
		result[i-1] = robj.Ptr.(ds.Sdshdr).Buf
	}
	core.AddReplyStatus(c, "[" + strings.Join(result, ",") + "]")
}