package core

// ZedisObject 是对特定类型的数据的包装
type ZedisObject struct {
	ObjectType int
	Ptr        interface{}
}

const ObjectTypeString = 1

// CreateObject 创建特定类型的object结构
func CreateObject(t int, ptr interface{}) (o *ZedisObject) {
	o = new(ZedisObject)
	o.ObjectType = t
	o.Ptr = ptr
	return
}
