package core

// Yedis的命令结构体，主要定义命令名称，命令处理函数和命令标志等
type YedisCommand struct {

	Name string //命令名称，如：set
	CommandProc YedisCommandProc //命令处理的函数
	Arity int //参数数量
	SFlags string //命令标识，标识是读写还是其他命令
	Flags int //命令的二进制标识，服务启动解析SFlags生成
	Calls int64 //从服务器启动到现在，命令执行的次数，用于统计
	Microseconds int64 //从服务器启动到现在，命令执行的总时间，Microseconds/Calls能计算出命令平均处理时间，用于统计
}

type YedisCommandProc func(c *YedisClients, s *YedisServer)