package core

import (
	"strconv"
)

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

//设置对象的编码方式来节省空间
//分配OBJ_ENCODING_INT，OBJ_ENCODING_RAW，OBJ_ENCODING_EMBSTR这三种编码方式
//此处简写，原版代码：https://github.com/antirez/redis/blob/30724986659c6845e9e48b601e36aa4f4bca3d30/src/object.c#L439
func TryObjectEncoding(lobj *YedisObject) *YedisObject {
	str := lobj.Ptr.(Sdshdr).Buf
	//判断str能否转int
	_, err := strconv.Atoi(str)
	if err == nil {
		lobj.Encoding = OBJ_ENCODING_INT
		return lobj
	}

	if len(str) <= 20 {
		lobj.Encoding = OBJ_ENCODING_EMBSTR
		return lobj
	}else {
		lobj.Encoding = OBJ_ENCODING_RAW
		return lobj
	}

}


func CreateSetObject() *YedisObject {
	ht := new(DictHt)
	ht.Table = make(DictMap, DEFAULT_HASH_LEN)
	ht.Size = DEFAULT_HASH_LEN
	ht.SizeMask = DEFAULT_HASH_LEN - 1
	ht.Used = 0
	o := CreateObject(REDIS_SET, OBJ_ENCODING_HT, ht)
	return o
}

//创建一个hash表对象
func CreateHashObject() *YedisObject {

	ht := new(DictHt)
	ht.Table = make(DictMap, DEFAULT_HASH_LEN)
	ht.Size = DEFAULT_HASH_LEN
	ht.SizeMask = DEFAULT_HASH_LEN - 1
	ht.Used = 0
	o := CreateObject(REDIS_HASH, OBJ_ENCODING_HT, ht)
	return o
}

//创建一个sds简单动态字符串对象
func CreateSdsObject(encodingType int, str string) *YedisObject {
	sdshdr := Sdsnew(str)
	return CreateObject(REDIS_STRING, encodingType, sdshdr)
}

//创建一个链表编码的Yedis对象
//代码：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/object.c#L217
func CreateLinkedListObject() *YedisObject {
	return CreateObject(REDIS_LIST, OBJ_ENCODING_LINKEDLIST, ListCreate())
}
