package sds

import (
	"Monica/go-yedis/core"
	"strconv"
)

//如果 key 已经存在并且是一个字符串， APPEND 命令将指定的 value 追加到该 key 原来值（value）的末尾
func AppendCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	robj := core.LookupKey(c.Db.Dict, c.Argv[1])
	if robj != nil {
		if sdshdr, ok := robj.Ptr.(core.Sdshdr); ok {
			//获取到字符串sdshdr对象,将Buf追加参数Argv[2]
			sdshdr.Buf = sdshdr.Buf + c.Argv[2].Ptr.(string)
			sdshdr.Len = sdshdr.Len + uint64(len(c.Argv[2].Ptr.(string)))
			robj.Ptr = sdshdr //TODO 直接覆盖原有Sds,通过指针修改不知道为什么没法成功
			s.Dirty++
			core.AddReplyStatus(c, strconv.FormatUint(sdshdr.Len, 10))
		}else {
			core.AddReplyStatus(c, "nil")
		}
	}else {
		core.AddReplyStatus(c, "nil")
	}
}