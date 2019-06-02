package znet

import (
	"fmt"
	"zinx/config"
	"zinx/ziface"
)

type MsgHandler struct {
	//存放路由集合的map
	Apis map[uint32]ziface.IRouter
	//取worker对应的消息队列 一个worker一个消息队列
	TaskQueue []chan ziface.IRequest
	//worker工作池的数量
	WorkerPoolSize uint32
}

//初始化方法
func NewMsgHandler() ziface.IMsgHandler {
	//给map开辟头空间
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, config.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: config.GlobalObject.WorkerPoolSize,
	}
}

//添加路由到map中
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	//判断新添加的msgID是否存在
	if _, OK := mh.Apis[msgId]; OK {
		fmt.Println("User is registered")
		return
	}
	//添加用户msgIDherouter的对应关系
	mh.Apis[msgId] = router
	fmt.Println("Added user ID succeeded", msgId)
}

//g根据msgID调度路由
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	//从request总获取用户Id
	router, OK := mh.Apis[request.GetMsg().GetMsgId()]
	if !OK {
		fmt.Println("Api ID=", request.GetMsg().GetMsgId(), " is not registered")
		return
	}
	//根据msdID 对不同的router进行调用
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)

}

//worker处理业务的goroutine函数
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("workerID", workerID, "is start...")
	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandler(req)

		}
	}
}

//启动worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	fmt.Println("worker pool start")
	//根据wokerpool数量 创建worker goroutine
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//给worker所绑定的消息对象开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, config.GlobalObject.MaxWorkerTaskLen)
		//启动一个worker阻塞等待的消息从对应的管道中进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//将消息添加到worker工作池中
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//平均分配消息给worker 确定当前消息让哪一个request来处理
	//根据workerID来分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//将request发送给对应的taskQueue
	mh.TaskQueue[workerID] <- request
}
