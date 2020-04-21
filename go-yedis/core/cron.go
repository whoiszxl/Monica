package core

import (
	"Monica/go-yedis/persistence"
	"Monica/go-yedis/utils"
	"fmt"
	"log"
)

//Redis的定时任务器，每秒钟调用config.hz次，默认是每秒十次
//其中Yedis实现的异步操作如下:
//1. 去Db.ExpireDict下清除过期的key
//2. 更新一些统计信息，比如说内存使用情况，内存最高占用，command的平均调用时间
//3. 调用bgsave备份数据到rdb，或者进行aof重写
//4. 清除超时客户端连接
func ServerCron(loop *AeEventLoop, server *YedisServer) int {

	//1. 更新服务器的时间
	updateCachedTime(server)

	//2. TODO 统计记录服务器执行命令的次数
	Run_with_period(1000, trackOperationsPerSecond, server)

	//3. TODO LRU实现

	//4. 记录服务器内存峰值
	RecordPeakMemory(server)

	//5. TODO 判断服务器是否收到SIGTERM信号，如果收到就关闭服务器

	//6. 客户端相关定时任务，关闭超时客户端
	clientCron()

	//7. 数据库相关定时任务
	databasesCron()

	//8. TODO 判断是否要执行aof重写, BGSAVE 和 BGREWRITEAOF 都没有在执行,有一个 BGREWRITEAOF 在等待的时候
	if server.RdbChildPid == -1 && server.AofChildPid == -1 && server.AofRewriteScheduled == 1 {
		persistence.RewriteAppendOnlyFileBackground()
	}

	//9. TODO 此时有BGSAVE或者BGREWRITEAOF在执行，需要接收完成信号来执行Handler，此处暂时先省略
	if server.RdbChildPid != -1 || server.AofChildPid != -1 {
		//接收信号
	}else {
		//如果没有后台重写aof和rdb bgsave, 检查是否需要执行bgsave
		//10. 检查更新的数量是否大于配置数量，还有时间是否超过了配置时间
		if server.Dirty >= server.SaveNumber && server.Unixtime - server.LastSaveTime > server.SaveNumber {
			persistence.RdbSaveBackground(server.RdbFileName)
		}
	}

	//11. TODO Trigger an AOF rewrite if needed, 不清楚为什么还要再一次触发rewriteAppendOnlyFileBackground()

	//12. 是否要将AOF缓冲区内容写入AOF文件中，因为AOF执行中，也有执行的命令，为了保持同步，需要将AOF缓冲区的数据也写入
	if server.AofFlushPostponedStart != 0 {
		flushAppendOnlyFile(0)
	}

	//13. 关闭需要异步关闭的客户端和暂停的客户端
	freeClientsInAsyncFreeQueue()
	clientsArePaused()

	//14. 增加任务调用计数器
	server.Cronloops++

	//15. 返回配置hz后，定时任务需要间隔的时间，hz默认为10，则每100毫秒执行一次当前serverCron函数
	return 1000/server.Hz
}


//统计记录服务器执行命令的次数
func trackOperationsPerSecond(server *YedisServer) int {
	fmt.Println("TODO 统计记录服务器执行命令的次数")
	return 1
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

//
func flushAppendOnlyFile(flags int) {

}

func freeClientsInAsyncFreeQueue() {

}

func clientsArePaused() {

}

func RecordPeakMemory(server *YedisServer) {
	currentUsedMemory := int64(utils.GetUsedMemory())
	if server.StatPeakMemory < currentUsedMemory {
		server.StatPeakMemory = currentUsedMemory
	}
}

//更新server的时间
func updateCachedTime(server *YedisServer) {
	millis := utils.CurrentTimeMillis()
	server.Unixtime = millis / 1e3
	server.Mstime = millis
}