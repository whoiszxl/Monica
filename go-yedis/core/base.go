package core

//通过key查询db中data，获取value值
func LookupKey(data Dict, userKey *YedisObject) (ret *YedisObject) {
	
	//TODO 每次都要遍历，后续找优化方法
	for key, val := range data {
		if key.Ptr.(Sdshdr).Buf == userKey.Ptr.(string) {
			return val
		}
	}
	return nil
}

func GetKeyObj(data Dict, userKey *YedisObject) (ret *YedisObject) {
	for key := range data {
		if key.Ptr.(Sdshdr).Buf == userKey.Ptr.(string) {
			return key
		}
	}
	return nil
}