package sds

import (
	"Monica/go-yedis/core"
	"strings"
)

//mget命令
func MgetCommand(c *core.YedisClients, s *core.YedisServer) {


	//循环获取
	var result = make([]string, c.Argc - 1)
	for i := 1; i < c.Argc; i++ {
		robj := core.LookupKey(c.Db.Dict, c.Argv[i])
		result[i-1] = robj.Ptr.(core.Sdshdr).Buf
	}
	core.AddReplyStatus(c, "[" + strings.Join(result, ",") + "]")
}