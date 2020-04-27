package sds

import (
	"Monica/go-yedis/core"
	"strconv"
)

//strlen命令
func StrlenCommand(c *core.YedisClients, s *core.YedisServer) {

	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		if sdshdr, ok := robj.Ptr.(core.Sdshdr); ok {
			core.AddReplyStatus(c, strconv.FormatUint(sdshdr.Len, 10))
		}else if intValue, err := strconv.Atoi(robj.Ptr.(string)); err == nil{
			core.AddReplyStatus(c, strconv.Itoa(intValue))
		}else {
			core.AddReplyStatus(c, "nil")
		}
	}else {
		core.AddReplyStatus(c, "nil")
	}
}