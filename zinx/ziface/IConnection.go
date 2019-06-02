package ziface
//链接模块
import (
	"net"
)

type IConnection interface {
	Start()
	Stop()
	//获取链接ID
	GetConnID() uint32
	//获取conn原生套接字
	GetTCPConnection() *net.TCPConn
	//获取远程客户IP端地址
	GetRemoteAddr() net.Addr
	//发送数据给客户端
	Send(msgId uint32,msgData []byte) error
	//设置属性
	SetProperty(key string, value interface{})
	//获取属性
	GetProperty(key string) (interface{}, error)
	//删除属性
	RemoveProperty(key string)
}

//定义业务处理接口
type HandleFunc func(request IRequest) error
