package string

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
	"strconv"
)

//strlen命令
func StrlenCommand(c *core.YedisClients, s *core.YedisServer) {

	robj := command.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		if sdshdr, ok := robj.Ptr.(ds.Sdshdr); ok {
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