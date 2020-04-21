package core

import (
	"Monica/go-yedis/persistence"
	"Monica/go-yedis/utils"
	"log"
)

//Redis的定时任务器，每秒钟调用config.hz次，默认是每秒十次
//其中Yedis实现的异步操作如下:
//1. 去Db.ExpireDict下清除过期的key
//2. 更新一些统计信息，比如说内存使用情况，内存最高占用，command的平均调用时间
//3. 调用bgsave备份数据到rdb，或者进行aof重写
//4. 清除超时客户端连接
func ServerCron(loop *AeEventLoop, server *YedisServer) int {

	//数据库相关定时任务
	databasesCron()

	//客户端相关定时任务
	clientCron()

	//记录服务器的内存峰值
	RecordPeakMemory()


	//如果没有后台重写aof和rdb bgsave, 检查是否需要执行bgsave
	//检查更新的数量是否大于配置数量，还有时间是否超过了配置时间
	if server.Dirty >= server.SaveNumber && server.Unixtime - server.LastSaveTime > server.SaveNumber {
		persistence.RdbSaveBackground(server.RdbFileName)
	}

	//TODO 判断是否要执行aof重写
	if server.RdbChildPid == -1 && server.AofChildPid == -1 {
		persistence.RewriteAppendOnlyFileBackground()
	}


	//增加任务调用计数器
	server.Cronloops++

	return 1000/server.Hz
}

//定义时间事件的接口，serverCron需要实现它
type PeriodProc func(server *YedisServer) int
//serverCron函数中每个任务都是每秒调用server.hz次，有些任务需要对调用次数进行限制，就需要用到这个方法
//描述：传入毫秒数，如果小于server.hz的执行间隔时间，便直接执行，如果大于，则判断当前执行的次数是否是hz的整数倍
//简单描述：就是传入多少毫秒就间隔多少毫秒执行
func Run_with_period(_ms_ int, proc PeriodProc, server *YedisServer) {
	if _ms_ <= 1000/server.Hz || server.Cronloops%((_ms_)/(1000/server.Hz)) == 0 {
		proc(server)
	}
}

//数据库相关处理的执行函数
func databasesCron() {
	log.Println("开始执行databasesCron", utils.CurrentTimeMillis())
}

//客户端相关定时任务
func clientCron() {

}

func RecordPeakMemory() {

}