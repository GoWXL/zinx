package ziface
//链接管理模块
type IConnmanger interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(connID uint32)
	//根据链接ID得到链接
	Get(connID uint32) (IConnection, error)
	//得到目前服务器链接总个数
	Len()int
	//清空全部链接的方法
	ClearConn()
}
