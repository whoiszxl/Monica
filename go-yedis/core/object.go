package core

import "Monica/go-yedis/ds"

// ZedisObject 是对特定类型的数据的包装
type YedisObject struct {
	ObjectType int
	Encoding   int
	RefCount   int         //
	Ptr        interface{} //Ptr存储了某一种数据结构
}

// CreateObject 创建特定类型的object结构
func CreateObject(objectType int, encodingType int, ptr interface{}) (o *YedisObject) {
	o = new(YedisObject)
	o.ObjectType = objectType
	o.Encoding = encodingType
	o.Ptr = ptr
	o.RefCount = 1
	//TODO LRU开发
	return
}

//设置对象的编码
func TryObjectEncoding(lobj *YedisObject) *YedisObject {

}


func CreateSdsObject(encodingType int, str string) *YedisObject {
	sdshdr := ds.Sdsnew(str)
	return CreateObject(REDIS_STRING, encodingType, sdshdr)
}

//创建一个链表编码的Yedis对象
//代码：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/object.c#L217
func CreateLinkedListObject() *YedisObject {
	return CreateObject(REDIS_LIST, OBJ_ENCODING_LINKEDLIST, ds.ListCreate())
}