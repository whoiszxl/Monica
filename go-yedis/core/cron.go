package core

import (
	"Monica/go-yedis/utils"
	"fmt"
)

const ACTIVE_EXPIRE_CYCLE_SLOW = 0
const ACTIVE_EXPIRE_CYCLE_FAST = 1
const REDIS_DBCRON_DBS_PER_CALL = 16
const ACTIVE_EXPIRE_CYCLE_FAST_DURATION = 1000
const ACTIVE_EXPIRE_CYCLE_SLOW_TIME_PERC = 25
const ACTIVE_EXPIRE_CYCLE_LOOKUPS_PER_LOOP = 20

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
	databasesCron(server)

	//8. TODO 判断是否要执行aof重写, BGSAVE 和 BGREWRITEAOF 都没有在执行,有一个 BGREWRITEAOF 在等待的时候
	if server.RdbChildPid == -1 && server.AofChildPid == -1 && server.AofRewriteScheduled == 1 {
		RewriteAppendOnlyFileBackground()
	}

	//9. TODO 此时有BGSAVE或者BGREWRITEAOF在执行，需要接收完成信号来执行Handler，此处暂时先省略
	if server.RdbChildPid != -1 || server.AofChildPid != -1 {
		//接收信号
	}else {
		//如果没有后台重写aof和rdb bgsave, 检查是否需要执行bgsave
		//10. 检查更新的数量是否大于配置数量，还有时间是否超过了配置时间
		if server.Dirty >= server.SaveNumber && server.Unixtime - server.LastSaveTime > server.SaveNumber {
			RdbSaveBackground(server.RdbFileName)
		}
	}

	//11. TODO Trigger an AOF rewrite if needed, 不清楚为什么还要再一次触发rewriteAppendOnlyFileBackground()

	//12. 是否要将AOF缓冲区内容写入AOF文件中，因为AOF执行中，也有执行的命令，为了保持同步，需要将AOF缓冲区的数据也写入
	server.AofFlushPostponedStart = server.Unixtime
	if server.AofFlushPostponedStart != 0 {
		flushAppendOnlyFile(server, 0)
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
func databasesCron(server *YedisServer) {
	// 如果是从服务器，不能主动清除过期键
	if server.ActiveExpireEnabled == 1 && server.Masterhost == "" {
		// 清除模式：CYCLE_SLOW ，尽量多清除过期键
		activeExpireCycle(server, ACTIVE_EXPIRE_CYCLE_SLOW)
	}
}

//删除数据库中过期的键
//校验的数据库数量不会超过REDIS_DBCRON_DBS_PER_CALL
//类型为ACTIVE_EXPIRE_CYCLE_FAST，执行快速模式，执行的时长不会超过EXPIRE_FAST_CYCLE_DURATION毫秒，并在EXPIRE_FAST_CYCLE_DURATION毫秒内不会再次执行
//类型为ACTIVE_EXPIRE_CYCLE_SLOW，执行正常模式，执行时间限制为REDIS_HS常量的一个百分比，百分比由 REDIS_EXPIRELOOKUPS_TIME_PERC 定义
func activeExpireCycle(server *YedisServer, expireType int) {

	current_db := 0 //当前执行的数据库编号
	timelimit_exit := 0 //是否达到了快速模式的时间限制
	last_fast_cycle := 0 //最后一个快速循环运行时

	j, iteration := 0, 0

	//每次处理的数据库数量
	dbs_per_call := REDIS_DBCRON_DBS_PER_CALL

	//函数开始执行的时间
	start := utils.CurrentTimeMillis()

	//快速模式
	if expireType == ACTIVE_EXPIRE_CYCLE_FAST {
		if  timelimit_exit == 0 {
			return
		}
		if start < last_fast_cycle + ACTIVE_EXPIRE_CYCLE_FAST_DURATION*2 {
			return
		}
		last_fast_cycle = start
	}

	//如果遇到了时间限制，这次需要对所有数据库进行扫描，避免过多过期键占用空间
	if dbs_per_call > server.DbNum || timelimit_exit != 0 {
		dbs_per_call = server.DbNum
	}

	//获取函数处理的微秒时间上限
	timelimit := 1000000 * ACTIVE_EXPIRE_CYCLE_SLOW_TIME_PERC/server.Hz/100
	timelimit_exit = 0
	if timelimit <= 0 {
		timelimit = 1
	}

	//快速模式最多运行FAST_DURATION微秒,默认值为1000微秒
	if expireType == ACTIVE_EXPIRE_CYCLE_FAST {
		timelimit = ACTIVE_EXPIRE_CYCLE_FAST_DURATION
	}

	//遍历数据库
	for j = 0; j < dbs_per_call; j++ {
		var expired int
		//获取需要处理的数据库，并将游标+1
		db := server.ServerDb[current_db % server.DbNum]
		current_db++

		for {
			var num,now,ttl_sum,ttl_samples int

			//获取当前库中有多少过期键值对
			num = len(db.Expires)
			if num == 0 {
				db.AvgTTL = 0
				break
			}
			now = utils.CurrentTimeMillis()

			//Redis是采用Dict字典，需要扩容等操作，Yedis不采用这种，没有复杂的查找操作，直接遍历压过去就好了
			//Redis Dict结构地址：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/dict.h#L135

			// 已处理过期键计数器,键总TTL计数,总共处理的键计数器
			expired,ttl_sum,ttl_samples = 0,0,0

			//每次最多处理20个键
			if num > ACTIVE_EXPIRE_CYCLE_LOOKUPS_PER_LOOP {
				num = ACTIVE_EXPIRE_CYCLE_LOOKUPS_PER_LOOP
			}

			//开始随机获取数据库中的带过期的键，判断是否真的过期了，过期删除并计数，没过期就接着循环
			//这里暂用一下，随机感觉比顺序遍历更消耗内存
			for num != 0 {
				key, value := DictGetRandomKey(db.Expires)
				ttl := value - now
				if activeExpireCycleTryExpire(db, key, value, now) == 1 {
					expired++
				}
				if ttl < 0 {
					ttl = 0
				}
				ttl_sum += ttl //键累计的TTL
				ttl_samples++ //键累计的个数
				num--
			}

			//统计平均TTL
			if ttl_samples > 0 {
				avg_ttl := ttl_sum / ttl_samples
				//第一次设置平均TTL
				if db.AvgTTL == 0 {
					db.AvgTTL = int64(avg_ttl)
				}
				//获取上次平均TTL和这次的平均值
				db.AvgTTL = (db.AvgTTL + int64(avg_ttl))/2
			}

			iteration++

			//每16次遍历执行一次,并且需要遍历的时间超过timelimit
			if iteration & 16 == 0 && (utils.CurrentTimeMillis() - start > timelimit){
				timelimit_exit = 1
			}

			if timelimit_exit == 1 {
				return
			}

			if expired < ACTIVE_EXPIRE_CYCLE_LOOKUPS_PER_LOOP/4 {
				break
			}
		}
	}

}

//判断是否过期，过期就删除这个key
func activeExpireCycleTryExpire(db *YedisDb, key *YedisObject, expiredTime int, now int) int {
	//已经过期
	if now > expiredTime {
		ret := db.Data[key]
		if ret == nil {
			return 0
		}
		//删除键
		DbDelete(db, key)

		//TODO 传播过期命令 propagateExpire
		//TODO 发送事件 notifyKeyspaceEvent
		//TODO 减少引用计数 decrRefCount
		return 1
	}else {
		return 0
	}
}

//客户端相关定时任务
func clientCron() {

}

//将server.aof_buf刷到文件中
//策略为everysec时，如果后台有fsync运行，可能会延迟flush操作
//force为1需要强制刷入，0则可能延迟
func flushAppendOnlyFile(server *YedisServer, force int) {

	//校验缓存中是否有值
	if server.AofBuf == "" {
		return
	}

	if server.AofFsync == AOF_FSYNC_EVERYSEC {
		//todo 判断是否有sync在后台运行 sync_in_progress = bioPendingJobsOfType(REDIS_BIO_AOF_FSYNC) != 0;
	}
	//TODO 推迟操作先不做

	server.AofFlushPostponedStart = 0


	//看C的实现，是先调用write写入aof文件，然后再判断是不是__linux__系统，linux系统再调fdatasync刷入，其他系统调fsync,这里简化一下吧，golang直接write就好了吧
	//Redis源代码地址：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884e75503fb8be1a322406f21fb455f67/src/aof.c#L441
	err := AppendToFile(server.AofFileName, server.AofBuf)
	if err != nil {
		fmt.Println("aof file write fail.")
		//TODO 写入出错，需要写入到日志里
		return
	}

	//更新写入后的AOF文件大小 TODO 更新的长度，不是文件大小
	server.AofCurrentSize += len(server.AofBuf)

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