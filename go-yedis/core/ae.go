package core

//事件处理器结构体
//Yedis服务器是事件驱动程序，事件类型有两种，分别为文件事件（socket读写事件，已经在yedis-server.go里写死了）和时间事件（就是定时任务）
// Server.Al<-AeEventLoop<-AeTimeEvent<-ServerCron
type AeEventLoop struct {

	Stop int //标识时间循环是否结束
	TimeEventHead AeTimeEvent //时间事件链表的头节点
}


//时间事件结构
//一些清除超时客户端连接，删除过期key，获取统计信息的定时任务会被封装
//多个时间事件会形成一个链表
//processTimeEvents处理逻辑：遍历时间事件列表，
type AeTimeEvent struct {

	Id int //时间事件唯一ID
	WhenSec int //时间事件触发的秒数
	WhenMs int //时间事件触发的毫秒数
	TimeProc *AeTimeProc //时间事件处理函数的指针
	FinalizerProc *AeEventFinalizerProc //删除时间事件节点之前会调用此函数
	clientData *YedisClients //对应的客户端对象的指针
	next *AeTimeEvent //指向下一个时间事件节点
}

type AeTimeProc func(server *YedisServer)

type AeEventFinalizerProc func(c *YedisClients, s *YedisServer)

//创建时间事件
func AeCreateTimeEvent(server *YedisServer, loop *AeEventLoop, milliseconds int, proc AeTimeProc, clients *YedisClients, finalizerProc *AeEventFinalizerProc) {



}