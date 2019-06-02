package ziface

//路由模块
type IRouter interface {
	//处理业务之前的方法
	PreHandle(request IRequest)
	//处理业务的主方法
	Handle(request IRequest)
	//处理业务的主方法
	PostHandle(request IRequest)
}
