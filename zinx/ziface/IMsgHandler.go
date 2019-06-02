package ziface

//多路由集合 抽象的消息管理模块
type IMsgHandler interface {
	//添加路由到map中
	AddRouter(msgId uint32, router IRouter)
	//g根据msgID调度路由
	DoMsgHandler(request IRequest)
	//启动worker工作池
	StartWorkerPool()
	//将消息添加到worker工作池中
	SendMsgToTaskQueue(request IRequest)
}
