package core

import "math/rand"

const DICT_OK = 0 // 操作成功
const DICT_ERR = 1 // 操作失败（或出错）

type YedisDb struct {
	ID int8 //数据库序号，默认Yedis有16个数据库，从0-15.
	Data Dict //Data中存储数据库中所有的键值对，Redis原命名是dict，这里采用data，感觉看着更舒服一点,源码地址https://github.com/antirez/redis/blob/30724986659c6845e9e48b601e36aa4f4bca3d30/src/server.h#L642
	Expires ExpireDict //存储键值对的过期时间
	AvgTTL int64 //数据库对象的平均TTL,用于统计
}

//使用Go原生数据结构map作为redis中dict结构体
type Dict map[*YedisObject]*YedisObject

//保存过期键值对的字典  map[键名]过期时间的时间戳
type ExpireDict map[*YedisObject]int

//获取一个随机的失效key
func DictGetRandomKey(dict ExpireDict) (*YedisObject, int){
	randomNumber := rand.Intn(len(dict))
	index := 0
	for key, value := range dict {
		if index == randomNumber {
			return key, value
		}
		index++
	}
	return nil, 0
}

//从数据库中删除给定的键值对和过期时间
func DbDelete(db *YedisDb, key *YedisObject) int {
	//删除键过期时间
	if len(db.Expires) > 0 {
		dictExpireDelete(db.Expires, key)
	}

	//删除键值对
	if dictDataDelete(db.Data, key) == DICT_OK {
		//TODO 集群模式要从slot中删除给定的键值对
		return 1
	}else {
		return 0
	}
}

//直接用go原生map直接删除一个key，美滋滋
//Redis实现很复杂，可以查看Redis代码：https://github.com/huangz1990/redis-3.0-annotated/blob/unstable/src/dict.c#L654
func dictExpireDelete(dict ExpireDict, key *YedisObject) int {
	delete(dict, key)
	return 1
}

func dictDataDelete(dict Dict, key *YedisObject) int {
	delete(dict, key)
	return 1
}


//将key添加到数据库中
//Redis源码：https://github.com/antirez/redis/blob/3.0/src/db.c#L93
func DbAdd(db *YedisDb, key string, value *YedisObject) {

}