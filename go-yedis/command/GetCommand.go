package command

import "Monica/go-yedis/core"

//get命令
func GetCommand(c *core.YedisClients, s *core.YedisServer) {

	robj := lookupKey(c.Db.Data, c.Argv[1])
	if robj != nil {
		core.AddReplyStatus(c, robj.Ptr.(string))
	}else {
		core.AddReplyStatus(c, "nil")
	}
}