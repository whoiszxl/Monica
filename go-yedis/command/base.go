package command

import (
	"Monica/go-yedis/core"
	"Monica/go-yedis/ds"
)

//通过key查询db中data，获取value值
func LookupKey(data core.Dict, userKey *core.YedisObject) (ret *core.YedisObject) {
	
	//TODO 每次都要遍历，后续找优化方法
	for key, val := range data {
		if key.Ptr.(ds.Sdshdr).Buf == userKey.Ptr.(string) {
			return val
		}
	}
	return nil
}

func GetKeyObj(data core.Dict, userKey *core.YedisObject) (ret *core.YedisObject) {
	for key := range data {
		if key.Ptr.(ds.Sdshdr).Buf == userKey.Ptr.(string) {
			return key
		}
	}
	return nil
}