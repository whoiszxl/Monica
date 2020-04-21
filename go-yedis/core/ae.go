package core

import (
	"Monica/go-yedis/utils"
)

//事件处理器结构体
//Yedis服务器是事件驱动程序，事件类型有两种，分别为文件事件（socket读写事件，已经在yedis-server.go里写死了）和时间事件（就是定时任务）
// Server.Al<-AeEventLoop<-AeTimeEvent<-ServerCron
type AeEventLoop struct {

	Stop int //时间处理器开关，标志是否结束
	TimeEventHead *AeTimeEvent //时间事件链表的头节点
	LastTime int //最后一次执行时间事件的时间
	Beforesleep AeBeforeSleepProc //事件执行前需要执行的函数
	TimeEventNextId int //用于生成时间事件 id
	Setsize int //已跟踪的文件描述符的最大数量，取得是最大客户端连接数加上一个REDIS_EVENTLOOP_FDSET_INCR，具体什么意思还没搞懂,暂时先在初始化里写死了
}


//时间事件结构
//一些清除超时客户端连接，删除过期key，获取统计信息的定时任务会被封装
//多个时间事件会形成一个链表
//processTimeEvents处理逻辑：遍历时间事件列表，
type AeTimeEvent struct {

	Id int //时间事件唯一ID
	WhenSec int //时间事件触发的秒数  存储
	WhenMs int //时间事件触发的毫秒数
	TimeProc AeTimeProc //时间事件处理函数的指针
	FinalizerProc *AeEventFinalizerProc //删除时间事件节点之前会调用此函数
	clientData *YedisClients //对应的客户端对象的指针
	next *AeTimeEvent //指向下一个时间事件节点
}

//定义时间事件的接口，serverCron需要实现它
type AeTimeProc func(loop *AeEventLoop, server *YedisServer) int

//定义事件释放时需要调用的接口  类似Spring中的切面
type AeEventFinalizerProc func( loop *AeEventLoop, clientData *YedisClients)

//定义事件之前需要执行的函数接口  类似Spring中的切面
type AeBeforeSleepProc func(s *YedisServer, loop *AeEventLoop)

//事件处理器的主循环执行函数
func AeMain(server *YedisServer) {
	//从服务器中获取时间事件循环对象，并把stop状态设置为关闭
	eventLoop := server.El
	eventLoop.Stop = 0

	for eventLoop.Stop == 0 {
		//判断是否需要在时间事件处理前执行before切面函数
		if eventLoop.Beforesleep != nil {
			eventLoop.Beforesleep(server, eventLoop)
		}

		//开始处理时间事件
		aeProcessEvents(server, eventLoop, AE_TIME_EVENTS)
	}
}

//处理所有时间事件，Redis中还需要处理文件事件，这里简化了就不调用文件事件了
//flags代表了事件类型
// 0:函数不处理，直接返回
// 1:AE_FILE_EVENTS,处理文件事件
// 2:AE_TIME_EVENTS,处理时间事件
// 1|2:AE_ALL_EVENTS,处理所有事件
// 4:AE_DONT_WAIT,处理完成所有非阻塞事件后立马返回
func aeProcessEvents(server *YedisServer, eventLoop *AeEventLoop, flags int) int {
	var processed = 0

	if flags != AE_TIME_EVENTS {
		return 0
	}

	//代码还需要对文件事件处理，此处不处理，Redis代码：https://github.com/huangz1990/redis-3.0-annotated/blob/8e60a75884/src/ae.c#L509

	//直接执行时间事件
	if flags == AE_TIME_EVENTS {
		processed += processTimeEvents(server, eventLoop)
	}

	return processed
}

//处理时间事件
func processTimeEvents(server *YedisServer, eventLoop *AeEventLoop) int {

	processed := 0
	te := new(AeTimeEvent)
	//maxId := 0
	now := utils.CurrentTimeMillis()

	//重置时间事件的运行时间防止因为时间穿插造成的事件处理混乱
	//如果当前时间小于上次执行时间事件的时间，说明穿插了，不能混在一个时间线上执行
	if now < eventLoop.LastTime {
		//取出头时间事件到中间变量
		te = eventLoop.TimeEventHead
		//将事件触发的描述置为0并指向下一个时间事件
		for te != nil {
			te.WhenSec = 0
			te = te.next
		}
	}

	//更新命令最后一次执行的时间
	eventLoop.LastTime = now

	//遍历链表执行那些whenSec时间到达了的事件
	te = eventLoop.TimeEventHead
	//maxId = eventLoop.TimeEventNextId - 1
	for te != nil {
		var id int

		//跳过无效事件，暂没搞明白什么意思
		//if te.Id > maxId {
		//	te = te.next
		//	continue
		//}

		//获取当前时间
		nowSec, nowMs := utils.CurrentSecondAndMillis()

		//如果当前时间等于大于事件的执行时间便执行
		if nowSec > te.WhenSec || (nowSec == te.WhenSec && nowMs >= te.WhenMs) {
			var retval int

			id = te.Id

			//执行时间事件处理器，并获取返回值
			retval = te.TimeProc(eventLoop, server)
			processed++

			//记录是否需要循环执行这个时间事件
			if retval != AE_NOMORE {
				//需要再次执行
				whenSec, whenMs := aeAddMillisecondsToNow(retval)
				te.WhenSec = whenSec
				te.WhenMs = whenMs
			}else {
				//不需要执行了,删除时间事件
				aeDeleteTimeEvent(eventLoop, id)
			}

			//将te继续放回表头，继续循环执行事件
			te = eventLoop.TimeEventHead
		}
		//else {
		//	te = te.next
		//}
	}

	return processed

}



//创建时间事件，并将时间事件保存到server全局对象中
func AeCreateTimeEvent(server *YedisServer, milliseconds int, proc AeTimeProc, clients *YedisClients, finalizerProc *AeEventFinalizerProc) int {

	//更新时间事件计数器
	server.El.TimeEventNextId++
	id := server.El.TimeEventNextId

	te := new(AeTimeEvent)
	te.Id = id

	// 设定处理事件的时间
	whenSec, whenMs := aeAddMillisecondsToNow(milliseconds)
	te.WhenSec = whenSec
	te.WhenMs = whenMs

	//设置事件处理器
	te.TimeProc = proc
	te.FinalizerProc = finalizerProc
	te.clientData = clients

	//将当前头部赋值给下一个，将当前的时间事件赋值给头部
	te.next = server.El.TimeEventHead
 	server.El.TimeEventHead = te
	return id
}

//在当前时间上加上 milliseconds 毫秒,并将添加了毫秒数的时间戳赋值到 sec 和 ms 中
func aeAddMillisecondsToNow(milliseconds int) (whenSec int, whenMs int) {

	//获取当前时间
	curSec, curMs := utils.CurrentSecondAndMillis()

	whenSec = curSec + milliseconds/1000
	whenMs = curMs + milliseconds%1000

	if whenMs >= 1000 {
		whenSec ++
		whenMs -= 1000
	}
	return whenSec, whenMs

}

//创建事件循环，给对象的属性赋值
func AeCreateEventLoop(setsize int) *AeEventLoop {
	// 创建事件状态对象
	eventLoop := new(AeEventLoop)

	//实际这里还要初始化文件事件，但是Yedis在事件循环中省略了，所以就没了。

	eventLoop.Setsize = setsize
	//初始化最近一次的执行时间
	eventLoop.LastTime = utils.CurrentTimeMillis()

	//初始化时间事件结构，后续在AeCreateTimeEvent函数中才将serverCron函数初始化进去
	eventLoop.TimeEventHead = nil
	eventLoop.TimeEventNextId = 0

	eventLoop.Stop = DISABLE
	eventLoop.Beforesleep = nil

	return eventLoop

}

// 寻找里目前时间最近的时间事件
// 因为链表是乱序的，所以查找复杂度为 O（N）
func aeSearchNearestTimer(eventLoop *AeEventLoop) *AeTimeEvent {
	te := eventLoop.TimeEventHead
	nearest := new(AeTimeEvent)


	for te.next != nil {
		//遍历事件链表，拿到事件最近的事件对象，其实里面就一个serverCron，Redis设计里也才两个时间事件，这样设计可能是为了后续开发扩展吧
		if te.WhenSec < nearest.WhenSec || (te.WhenSec == nearest.WhenSec && te.WhenMs < nearest.WhenMs) {
			nearest = te
		}
		te = te.next
	}

	return nearest
}


//停止事件处理器
func aeStop(eventLoop *AeEventLoop) {
	eventLoop.Stop = ENABLE
}

//删除时间事件，这个用不到，因为Yedis中暂时还没有单次定时事件，都是循环时间事件
func aeDeleteTimeEvent(loop *AeEventLoop, id int) {

}