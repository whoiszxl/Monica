package entity

import "Monica/go-yedis/core"

// ZedisObject 是对特定类型的数据的包装
type YedisObject struct {
	ObjectType int
	Encoding int
	RefCount int //
	Ptr interface{} //Ptr存储了某一种数据结构
}

// CreateObject 创建特定类型的object结构
func CreateObject(objectType int, ptr interface{}) (o *YedisObject) {
	o = new(YedisObject)
	o.ObjectType = objectType
	o.Encoding = core.OBJ_ENCODING_RAW
	o.Ptr = ptr
	o.RefCount = 1
	//TODO LRU开发
	return
}
