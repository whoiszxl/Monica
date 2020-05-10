package set

import (
	"Monica/go-yedis/core"
)

//判断元素是否属于当前set
func SismemberCommand(c *core.YedisClients, s *core.YedisServer) {
	//搜索key是否存在数据库中
	set := core.LookupKey(c.Db.Dict, c.Argv[1])
	if set != nil {
		htTable := set.Ptr.(*core.DictHt).Table
		//遍历判断输入的参数是否属于htTable,暂时先遍历查找 entry := htTable[c.Argv[2].Ptr.(core.Sdshdr)]

		for k, _ := range htTable {
			if k.(core.Sdshdr).Buf == c.Argv[2].Ptr.(core.Sdshdr).Buf {
				core.AddReplyStatus(c, "(integer) 1")
				return
			}
		}
	}
	core.AddReplyStatus(c, "(integer) 0")
	return
}