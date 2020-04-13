package main

import "fmt"
import "flag"
import "os"
import "os/signal"
import "syscall"
import "net"
import "log"
import "time"
import "Monica/go-redis/core"
import "Monica/go-redis/persistence"

const (
	DefaultAofFile = "./zedis.aof"
)

//创建服务端实例
var zedis = new(core.Server)

func main() {

	//flag.Parse()

	//获取输入命令的个数，如为1，则为查询版本使用指南等基础操作、
	//os.Args 参数0：文件全路径名 参数1：实际参数 ...
	cmdLen := len(os.Args)
	cmdArgs := os.Args

	//无命令参数输出帮助
	if cmdLen < 2 {
		help()
		os.Exit(1)
	}

	//单个参数输出基础帮助信息
	if cmdLen == 2 {
		baseHelp(cmdArgs[1])
	}

	//读取命令行输入的ip和端口
	var ip = flag.String("ip", "localhost", "redis服务端IP")
	var port = flag.String("port", "6380", "redis服务端PORT")
	flag.Parse()

	
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go signalHandler(c)

	//初始化服务端实例
	initServer()

	//连接网络
	netListen, err := net.Listen("tcp", *ip + ":" + *port)
	if err != nil {
		log.Print("net listen err")
	}

	//延迟关闭网络连接
	defer netListen.Close()


	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		go handle(conn)
	}
}

//初始化服务端实例
func initServer() {
	zedis.Pid = os.Getpid() //获取进程ID
	zedis.DbNum = 16 //配置db数量
	initDb() //初始化db
	zedis.Start = time.Now().UnixNano() / 1000000 //记录开始运行时间

	zedis.AofFilename = DefaultAofFile

	getCommand := &core.ZedisCommand{Name: "get", Proc: core.GetCommand}
	setCommand := &core.ZedisCommand{Name: "set", Proc: core.SetCommand}

	zedis.Commands = map[string]*core.ZedisCommand{
		"get": getCommand,
		"set": setCommand,
	}

	LoadData()
}

//初始化DB
func initDb() {
	zedis.Db = make([]*core.ZedisDb, zedis.DbNum)
	for i:=0; i<zedis.DbNum; i++ {
		zedis.Db[i] = new(core.ZedisDb)
		zedis.Db[i].Dict = make(map[string]*core.ZedisObject, 100)
	}
}

// 处理请求
func handle(conn net.Conn) {
	c := zedis.CreateClient()
	for {
		err := c.ReadQueryFromClient(conn)

		if err != nil {
			log.Println("readQueryFromClient err", err)
			return
		}
		err = c.ProcessInputBuffer()
		if err != nil {
			log.Println("ProcessInputBuffer err", err)
			return
		}
		
		zedis.ProcessCommand(c)
		response2Client(conn, c)
	}
}

// 响应返回给客户端
func response2Client(conn net.Conn, c *core.Client) {
	conn.Write([]byte(c.Buf))
}

// 读取客户端请求信息
func readQueryFromClient(conn net.Conn) (buf string, err error) {
	buff := make([]byte, 512)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("conn.Read err!=nil", err, "---len---", n, conn)
		conn.Close()
		return "", err
	}
	buf = string(buff)
	return buf, nil
}

func signalHandler(c chan os.Signal) {
	for s := range c {
		switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				ExitFunc()
			default:
				fmt.Println("other", s)
		}
	}
}


func LoadData() {
	fmt.Println("开始加载磁盘数据到内存")
	c := zedis.CreateClient()
	pros := persistence.ReadAof(zedis.AofFilename)
	for _, v := range pros {
		c.QueryBuf = string(v)
		err := c.ProcessInputBuffer()
		if err != nil {
			log.Println("ProcessInputBuffer err", err)
		}
		zedis.ProcessCommand(c)
	}
	fmt.Println("加载磁盘数据到内存结束")
}

func ExitFunc()  {
    fmt.Println("开始退出...")
    fmt.Println("执行清理...")
    fmt.Println("结束退出...")
    os.Exit(0)
}

//处理基础操作，查询版本，使用指南等
func baseHelp(command string) {
	switch command {
		case "-v":
			version()
		case "-version":
			version()
		case "version":
			version()
		case "-h":
			help()
		case "-help":
			help()
		case "help":
			help()
		default:
			printError()
	}
}

func version() {
	fmt.Println("Redis server v=3.2.12 sha=00000000:0 malloc=jemalloc-3.6.0 bits=64 build=7897e7d0e13773f")
	os.Exit(0)
}


func help() {
	fmt.Println("Usage: ./zedis-server [/path/to/redis.conf] [options]")
	fmt.Println("       ./zedis-server - (read config from stdin)")
	fmt.Println("       ./zedis-server -v or --version")
	fmt.Println("       ./zedis-server -h or --help")
	fmt.Println("       ./redis-server --test-memory <megabytes>")
	fmt.Println("Examples:")
	fmt.Println("       ./zedis-server (run the server with default conf)")
	fmt.Println("       ./zedis-server /etc/redis/6380.conf")
	fmt.Println("       ./zedis-server --port 7777")
	fmt.Println("       ./zedis-server --port 7777 --slaveof 127.0.0.1 8888")
	fmt.Println("       ./zedis-server /etc/myredis.conf --loglevel verbose")
	fmt.Println("Sentinel mode:")
	fmt.Println("       ./zedis-server /etc/sentinel.conf --sentinel")
	os.Exit(0)
}

func printError() {
	fmt.Println("命令输入错误")
	fmt.Println(os.Getpid())
}