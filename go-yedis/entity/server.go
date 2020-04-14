package entity

//服务端结构体
//结构体存储Yedis服务器的所有信息，包括但不限于数据库，配置参数,
//命令表，监听端口地址，客户端列表，RDB,AOF持久化信息等
//RedisServer原结构体地址：https://github.com/antirez/redis/blob/30724986659c6845e9e48b601e36aa4f4bca3d30/src/server.h#L1027
type YedisServer struct {

	/* 基础配置 */

	Pid int //主进程的PID编号
	ConfigFile string //配置文件绝对路径
	DbNum int //数据库的数量，可以通过yedis.conf配置，默认16个
	yedisDb []*YedisDb //储存数据库的数组

	//serverCron函数执行频率,最小值1，最大值500，Redis-3.0.0默认是10，代表每秒执行十次serverCron函数
	//serverCron函数执行类似清除过期键，处理超时连接等任务
	//Redis实际还有dynamic_hz和config_hz，分别是根据客户端数量动态调整的和配置调整的
	Hz int

	//命令字典
	//key：字符串类型命令，如: set get ttl等
	//value：命令的实际操作，YedisCommand的指针
	Commands map[string]*YedisCommand

	/* 网络配置 */
	BindAddr string //绑定运行的IP地址，简化为1个，Redis有多个
	Port int32 // Yedis服务器监听的端口号，可以通过yedis.conf配置，默认端口6380
	NextClientId int64 //下一个客户端的唯一ID
	Clients map[string]*YedisClients //当前连接的可用客户端
	ClientsToClose map[string]*YedisClients //当前关闭的客户端

	/* RDB persistence持久化 */
	Dirty int64 //存储上次数据变动前的长度
	RdbFileName string //rdb文件名
	RdbCompression int //是否对rdb使用压缩
	LastSaveTime int64 //最后一次保存的时间


	/* AOF persistence持久化 */
	AofEnabled int //是否开启Aof
	AofState int //aof状态，[0: OFF] [1: ON] [2: WAIT_REWRITE]
	AofFileName string //aof文件名
	AofCurrentSize int //aof文件当前大小
	AofBuf []string //aof缓冲区，在进入事件循环前写入


	/* 仅用于统计使用的字段，仅取部分 */
	StatStartTime int64 //服务启动时间
	StatNumCommands int16 //命令数量
	StatNumConnections int16 //连接数量
	StatExpiredKeys int64 //失效key的数量


	/* 系统硬件信息 */
	SystemMemorySize int64  //系统内存大小

}


