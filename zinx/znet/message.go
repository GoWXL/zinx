package znet

import "zinx/ziface"

//消息模块
type Message struct {
	//数据ID
	Id uint32
	//数据内容
	Data []byte
	//数据长度
	DataLen uint32
}

//初始化NewMessage方法
func NewMessage(id uint32, data []byte) ziface.IMessage {
	return &Message{
		Id:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

//得到数据ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//得到数据长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

//得到数据内容
func (m *Message) GetMsgData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id=id
}
func (m *Message) SetData(data []byte) {
	m.Data=data
}
func (m *Message) SetDatalen(len uint32) {
	m.DataLen=len
}
