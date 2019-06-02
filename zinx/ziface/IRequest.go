package ziface
//链接和数据绑定模块
type IRequest interface {
	//得到当前请求的链接
	GetConnection() IConnection
	/*//得到链接的数据
	GetData() []byte
	//得到数据的长度
	GetDataLen()int*/
	//得到请求的消息
	GetMsg() IMessage

}
