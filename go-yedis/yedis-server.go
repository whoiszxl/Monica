package main

import (
	"Monica/go-yedis/command"
	"Monica/go-yedis/command/sds"
	"Monica/go-yedis/core"
	"Monica/go-yedis/utils"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	//默认配置文件路径
	defaultConfigPath = "yedis.conf"

	//默认数据库的键值对初始容量
	defaultDbDictCapacity = 100
)

//创建服务端实例
var yedis = new(core.YedisServer)

func main() {

	//获取用户输入的参数
	cmdArgs := os.Args

	//第一个参数是文件路径， 读取配置
	var configPath = defaultConfigPath
	if len(cmdArgs) > 1 && strings.LastIndex(cmdArgs[1], ".conf") != -1 {
		fmt.Println("读取配置")
		configPath = cmdArgs[1]
	} else if len(cmdArgs) == 2 {
		utils.BaseHelp(cmdArgs[1])
	}

	//获取配置
	netConfig, dbConfig, aofConfig := utils.ReadConfig(configPath)
	fmt.Println("初始化yedis.conf网络参数", netConfig)
	fmt.Println("初始化yedis.conf数据库参数", dbConfig)
	fmt.Println("初始化yedis.confAOF持久化参数", aofConfig)

	//读取命令行输入的ip和端口,并将命令行获取的值回写
	var netBind = flag.String("ip", netConfig.NetBind, "redis服务端IP")
	var netPort = flag.String("port", netConfig.NetPort, "redis服务端PORT")
	flag.Parse()
	netConfig.NetBind = *netBind
	netConfig.NetPort = *netPort

	host := *netBind + ":" + *netPort
	log.Println("Redis实例化地址：" + host)

	//监听退出事件做相应处理
	utils.ExitHandler()

	//初始化服务端实例
	initServer(netConfig, dbConfig, aofConfig, configPath)

	// 事件循环有点复杂了，先直接go xx试下

	core.AeCreateTimeEvent(yedis, yedis.El, 1, core.ServerCron, nil, nil)


	//初始化网络监听并延时关闭
	//redis3.0代码：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/redis.c#L1981
	netListener, err := net.Listen("tcp", host)
	if err != nil {
		log.Println("net listen err:", err)
	}
	defer netListener.Close()

	//循环监听新连接，将新连接放入go协程中处理
	for {
		conn, err := netListener.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

//处理连接请求
func handle(conn net.Conn) {
	//通过服务器给新的请求创建一个连接
	c := yedis.CreateClient()
	for {
		//从连接中读取命令，并写入到Client对象中
		err := c.ReadCommandFromClient(conn)
		if err != nil {
			log.Println("ReadCommandFromClient err", err)
			return
		}

		//解析命令到Client的Argv中
		err = c.ProcessCommandInfo()
		if err != nil {
			log.Println("ProcessCommandInfo err", err)
			continue
		}

		//执行命令
		yedis.ExecuteCommand(c)

		//响应客户端
		response2Client(conn, c)

	}
}

// 响应返回给客户端
func response2Client(conn net.Conn, c *core.YedisClients) {
	_, err := conn.Write([]byte(c.Reply))
	//log.Println("响应的response字节数为：", responseSize)
	utils.ErrorVerify("消息响应客户端失败", err, false)
}

//初始化服务端实例, 将yedis.conf配置写入server实例
//redis3.0代码地址：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/redis.c#L3952
func initServer(netConfig utils.NetConfig, dbConfig utils.DbConfig, aofConfig utils.AofConfig, configPath string) {
	//1. 写入基础配置
	yedis.Pid = os.Getpid()            //获取进程ID
	yedis.ConfigFile = configPath      //配置文件绝对路径
	yedis.DbNum = dbConfig.DbDatabases //配置db数量
	yedis.Hz = dbConfig.Hz             //配置任务执行频率
	initDb()                           //初始化server中的16个数据库

	//2. 网络配置
	yedis.BindAddr = netConfig.NetBind //配置绑定IP地址
	yedis.Port = netConfig.NetPort     //配置端口号

	//3. RDB persistence持久化
	yedis.Dirty = 1                           //存储上次数据变动前的长度
	yedis.RdbFileName = dbConfig.DbDbfilename //rdb文件名
	yedis.RdbCompression = core.DISABLE       //TODO 是否对rdb使用压缩
	yedis.SaveTime = dbConfig.DbSavetime      //指定在多长时间内，有多少次更新操作，就将数据同步到数据文件，默认：300秒内10次更新操作就同步数据到文件
	yedis.SaveNumber = dbConfig.DbSavenumber  //

	//4. AOF persistence持久化
	if aofConfig.AofAppendonly == "no" { //配置是否开启aof：number
		yedis.AofEnabled = 0
	} else {
		yedis.AofEnabled = 1
	}
	yedis.AofState = aofConfig.AofAppendonly        //配置是否开启aof：字符串
	yedis.AofFileName = aofConfig.AofAppendfilename //配置aof文件名
	yedis.AofSync = aofConfig.AofAppendfsync        //配置同步文件的策略

	//5. 仅用于统计使用的字段，仅取部分
	yedis.StatStartTime = time.Now().UnixNano() / 1000000 //记录服务启动时间
	yedis.StatNumCommands = int16(len(yedis.Commands))    //支持的命令数量
	yedis.StatNumConnections = int16(0)                   //当前连接数量
	yedis.StatExpiredKeys = int64(0)                      //当前失效key的数量

	//6. 系统硬件信息
	memInfo, err := mem.VirtualMemory() //获取机器内存信息
	utils.ErrorVerify("获取机器内存信息失败", err, true)
	yedis.SystemAllMemorySize = memInfo.Total     //机器总内存大小 单位：b
	yedis.SystemAvailableSize = memInfo.Available //机器可用内存大小 单位：b
	yedis.SystemUsedSize = memInfo.Used           //机器已用内存大小 单位：b
	yedis.SystemUsedPercent = memInfo.UsedPercent //机器已用内存百分比

	percent, err := cpu.Percent(time.Second, false)
	utils.ErrorVerify("获取机器CPU信息失败", err, true)
	yedis.SystemCpuPercent = percent[0] //CPU使用百分比情况

	//初始化服务支持命令
	getCommand := &core.YedisCommand{Name: "get", CommandProc: sds.GetCommand}
	setCommand := &core.YedisCommand{Name: "set", CommandProc: sds.SetCommand}
	strlenCommand := &core.YedisCommand{Name: "strlen", CommandProc: sds.StrlenCommand}
	appendCommand := &core.YedisCommand{Name: "append", CommandProc: sds.AppendCommand}
	getrangeCommand := &core.YedisCommand{Name: "getrange", CommandProc: sds.GetrangeCommand}
	mgetCommand := &core.YedisCommand{Name: "mget", CommandProc: sds.MgetCommand}

	incrCommand := &core.YedisCommand{Name: "incr", CommandProc: sds.IncrCommand}
	incrbyCommand := &core.YedisCommand{Name: "incrby", CommandProc: sds.IncrbyCommand}
	decrCommand := &core.YedisCommand{Name: "decr", CommandProc: sds.DecrCommand}
	decrbyCommand := &core.YedisCommand{Name: "decrby", CommandProc: sds.DecrbyCommand}

	pexpireatCommand := &core.YedisCommand{Name: "pexpireat", CommandProc: sds.PexpireatCommand}

	infoCommand := &core.YedisCommand{Name: "info", CommandProc: command.InfoCommand}

	yedis.Commands = map[string]*core.YedisCommand{
		"get":      getCommand,
		"set":      setCommand,
		"strlen":   strlenCommand,
		"append":   appendCommand,
		"getrange": getrangeCommand,
		"mget":     mgetCommand,
		"info":     infoCommand,

		"incr":   incrCommand,
		"incrby": incrbyCommand,
		"decr":   decrCommand,
		"decrby": decrbyCommand,

		"pexpireat": pexpireatCommand,
	}

}

//初始化数据库
func initDb() {
	//创建一个储存数据库对象的切片
	yedis.ServerDb = make([]*core.YedisDb, yedis.DbNum)
	for i := 0; i < yedis.DbNum; i++ {
		//创建YedisDb数据库对象并对其中数据库ID和键值对字段赋值
		//键值对容量暂写死为200
		yedis.ServerDb[i] = new(core.YedisDb)
		yedis.ServerDb[i].ID = int8(i)
		yedis.ServerDb[i].Data = make(core.Dict, defaultDbDictCapacity)
		yedis.ServerDb[i].Expires = make(core.ExpireDict, defaultDbDictCapacity)
		yedis.ServerDb[i].AvgTTL = 0

	}
}
