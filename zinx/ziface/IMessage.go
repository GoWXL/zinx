package ziface

//消息模块
type IMessage interface {
	//得到数据ID
	GetMsgId() uint32
	//得到数据长度
	GetMsgLen() uint32
	//得到数据内容
	GetMsgData() []byte

	SetMsgId(uint32)
	SetData([]byte)
	SetDatalen(uint32)
}
