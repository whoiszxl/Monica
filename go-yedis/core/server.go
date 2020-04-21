package core

//服务端结构体
//结构体存储Yedis服务器的所有信息，包括但不限于数据库，配置参数,
//命令表，监听端口地址，客户端列表，RDB,AOF持久化信息等
//RedisServer原结构体地址：https://github.com/antirez/redis/blob/30724986659c6845e9e48b601e36aa4f4bca3d30/src/server.h#L1027
type YedisServer struct {

	/* 基础配置 */

	Pid int //主进程的PID编号
	ConfigFile string //配置文件绝对路径
	DbNum int //数据库的数量，可以通过yedis.conf配置，默认16个
	ServerDb []*YedisDb //储存数据库的数组
	Unixtime int //每一个cron定时任务都会更新的时间
	Mstime int //和unixtime一样，只是这个是毫秒
	El *AeEventLoop //所有的事件，链表结构，一般只有serverCron的事件
	Cronloops int //命令执行次数的计数器
	ShutdownAsap int //关闭服务器的标识，1：需要关闭  0：不关闭
	Requirepass string //请求时需要验证的密码，不设置则不校验

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
	Port string // Yedis服务器监听的端口号，可以通过yedis.conf配置，默认端口6380
	NextClientId int64 //下一个客户端的唯一ID
	MaxClients int //最大的客户端连接数
	Clients map[string]*YedisClients //当前连接的可用客户端
	ClientsToClose map[string]*YedisClients //当前关闭的客户端

	/* RDB persistence持久化 */
	Loading int //是否是正在载入中的状态 1：载入中 0：载入完成
	RdbChildPid int //执行bgsave子进程的pid，默认未执行状态为-1
	Dirty int //存储上次数据变动前的长度
	RdbFileName string //rdb文件名
	RdbCompression int //是否对rdb使用压缩
	LastSaveTime int //最后一次保存RDB的时间
	SaveTime int //rdb bgsave存储策略 时间  多少秒内有多少次修改则执行bgsave，原版是用数组保存，可以支持多个策略，此处先简写
	SaveNumber int //rdb bgsave存储策略 次数


	/* AOF persistence持久化 */
	AofChildPid int //执行aof重写的子进程id，默认未执行状态为-1
	AofEnabled int //是否开启Aof
	AofState string //aof状态，[0: OFF] [1: ON] [2: WAIT_REWRITE]
	AofFileName string //aof文件名
	AofCurrentSize int //aof文件当前大小
	AofBuf []string //aof缓冲区，在进入事件循环前写入
	AofSync string //更新模式：everysec: 每秒同步一次（折中，默认值，多用此配） no：表示等操作系统进行数据缓存同步到磁盘(效率高，不安全)  always：表示每次更新操作后手动调用fsync()将数据写到磁盘（效率低，安全，一般不采用）
	AofRewriteMinSize int //aof执行aof重新的最小大小
	AofRewriteScheduled int //AOF是否在执行重写，重写的时候需要阻塞其他aof和rdb任务，在bgrewriteaofCommand执行的时候需要将它设置为1，在success handler中需要设置回0
	AofFlushPostponedStart int //存储unix时间，推迟write flush的时间

	/* 仅用于统计使用的字段，仅取部分 */
	StatStartTime int64 //服务启动时间
	StatNumCommands int16 //命令数量
	StatNumConnections int16 //连接数量
	StatExpiredKeys int64 //失效key的数量
	StatPeakMemory int64 //服务器内存的峰值


	/* 系统硬件信息 */
	SystemAllMemorySize uint64  //系统内存大小
	SystemAvailableSize uint64 //系统可用内存
	SystemUsedSize uint64 //系统已用内存
	SystemUsedPercent float64 //内存使用百分比
	SystemCpuPercent float64 //CPU使用百分比

}


