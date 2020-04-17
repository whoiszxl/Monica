package core

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


//将key添加到数据库中
//Redis源码：https://github.com/antirez/redis/blob/3.0/src/db.c#L93
func DbAdd(db *YedisDb, key string, value *YedisObject) {

}