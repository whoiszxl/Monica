package core

import "math/rand"

const DICT_OK = 0 // 操作成功
const DICT_ERR = 1 // 操作失败（或出错）

type YedisDb struct {
	ID           int8       //数据库序号，默认Yedis有16个数据库，从0-15.
	Dict         Dict       //Data中存储数据库中所有的键值对，Redis原命名是dict，这里采用data，感觉看着更舒服一点,源码地址https://github.com/antirez/redis/blob/30724986659c6845e9e48b601e36aa4f4bca3d30/src/server.h#L642
	Expires      ExpireDict //存储键值对的过期时间
	AvgTTL       int64      //数据库对象的平均TTL,用于统计
	BlockingKeys Dict       //阻塞状态的所有键
	ReadyKeys    Dict       //等待解除阻塞的键
}

//使用Go原生数据结构map作为redis中dict结构体，实际Redis使用的是自己构建的一个Dict结构
//C语言因为没有自带的字典，使用的是数组，通过将输入的键通过hash计算（dictGenHashFunction）得到一串大数字，再将大数字通过
//和数组容量进行取余得出，因为这样操作会造成hash冲突，所以储存的值对象还会有一个next成员指针变量，形成一个个单向链表。
//查找的时候则先对键进行hash取余，取出单向链表后再一个个与key进行比对
//Redis源码：https://github.com/antirez/redis/blob/4d4c8c8a40/src/dict.h#L76
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
	if dictDataDelete(db.Dict, key) == DICT_OK {
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
func DbAdd(db *YedisDb, key *YedisObject, value *YedisObject) {
	db.Dict[key] = value
}

func SelectDb(c *YedisClients, s *YedisServer, id int) int {
	if id < 0 || id >= s.DbNum {
		return REDIS_ERR
	}

	//切换数据库
	c.Db = s.ServerDb[id]
	return REDIS_OK
}