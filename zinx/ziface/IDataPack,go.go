package ziface

//拆包封包模块 解决粘包问题
type IDataPack interface {
	//获取头部长度
	GetHeadLen() uint32
	//封包
	Pack() ([]byte, error)
	//拆包
	UnPack() (IMessage, error)
}
