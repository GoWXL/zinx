package znet

import "zinx/ziface"

//链接和数据绑定模块
type Request struct {
	//链接信息
	conn ziface.IConnection
	/*//数据内容
	data []byte
	//数据长度
	len int*/
	//得到消息的数据
	msg ziface.IMessage
}

//初始化方法
func NewRequest(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	req := &Request{
		conn: conn,
		msg:  msg,
	}

	return req
}

//得到当前请求的链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

/*//得到链接的数据
func (r *Request) GetData() []byte {
	return r.data
}

//得到数据的长度
func (r *Request) GetDataLen() int {
	return r.len
}*/
//得到请求的消息
func (r *Request) GetMsg() ziface.IMessage {
	return r.msg
}
