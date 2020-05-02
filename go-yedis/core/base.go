package core

//通过key查询db中data，获取value值
func LookupKey(data Dict, userKey *YedisObject) (ret *YedisObject) {
	
	//TODO 需要优化成hash取余方式，作O（n）复杂度
	for key, val := range data {
		if key.Ptr.(Sdshdr).Buf == userKey.Ptr.(Sdshdr).Buf {
			return val
		}
	}
	return nil
}

func GetKeyObj(data Dict, userKey *YedisObject) (ret *YedisObject) {
	for key := range data {
		if key.Ptr.(Sdshdr).Buf == userKey.Ptr.(Sdshdr).Buf {
			return key
		}
	}
	return nil
}