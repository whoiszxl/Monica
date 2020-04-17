package core

//
//RedisClient原结构体地址：https://github.com/antirez/redis/blob/30724986659c6845e9e48b601e36aa4f4bca3d30/src/server.h#L765
type YedisClients struct {
	Name string //客户端名称
	Argc int // 当前执行命令的参数的个数
	Argv []*YedisObject //当前执行命令的参数
	Db *YedisDb //指向当前选择数据库的指针
	QueryBuf string //积累客户端查询的缓冲区, 暂用string，后更新用sds
	Reply string //需要发送回客户端的回复信息
	Cmd *YedisCommand //待执行的客户端命令
	LastCommand *YedisCommand //上一个执行的
}