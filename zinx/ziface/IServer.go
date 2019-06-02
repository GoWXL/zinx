package ziface
//服务器模块
type IServer interface {
	Start()
	Stop()
	Server()
	//添加路由功能
	AddRouter(msgId uint32, router IRouter)
	GetConnMgr() IConnmanger
	//注册创建链接之后调用的hook函数的方法
	AddOnConnStart(hookFunc func(conn IConnection))
	//注册销毁链接之前调用的hook函数的方法
	AddOnConnStop(hookFunc func(conn IConnection))
	//调用创建链接之后调用的hook函数的方法
	CallOnConnStart(conn IConnection)
	//调用函数销毁之前调用的hook函数的方法
	CallOnConnStop(conn IConnection)
}
