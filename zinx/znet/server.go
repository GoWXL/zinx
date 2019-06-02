package znet

import (
	"fmt"
	"net"
	"zinx/config"
	"zinx/ziface"
)
//服务器模块
type Server struct {
	//服务器IP
	IPVersion string
	IP        string
	//服务器PORT
	Port int
	//服务器名字
	Name string
	//Router ziface.IRouter
	MsgHandler ziface.IMsgHandler
	//添加链接管理模块
	connManger ziface.IConnmanger
	//该server创建之后自动调用的hook函数指针
	OnConnStart func(conn ziface.IConnection)
	//该server链接销毁之前自动调用hook函数指针
	OnConnStop func(conn ziface.IConnection)
}

/*//定义一个回显业务
func CallBackBusi(request ziface.IRequest) error  {
	fmt.Println("【conn Handle】 CallBack..")
	c := request.GetConnection().GetTCPConnection()
	buf := request.GetData()
	cnt := request.GetDataLen()
	_, err := c.Write(buf[:cnt])
	if err != nil {
		fmt.Println("write back err", err)
		return err
	}
	return nil
}*/

//初始化new方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      config.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        config.GlobalObject.Host,
		Port:      config.GlobalObject.Port,
		//Router:    nil,
		MsgHandler: NewMsgHandler(),
		connManger: NewConnManger(),
	}
	return s
}
func (s *Server) Start() {
	fmt.Printf("[start] server Listenner at IP :%s,Port :%d,is starting..\n", s.IP, s.Port)
	s.MsgHandler.StartWorkerPool()
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error:", err)
		return
	}
	listenner, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listn", s.IPVersion, "err", err)
		return
	}
	var cid uint32
	cid = 0
	go func() {
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accepterr", err)
				continue
			}
			//判断当前server链接数量是否已经最大值
			if s.connManger.Len() >= int(config.GlobalObject.MaxConn) {
				//当前链接已经满了
				fmt.Println("---> Too many Connection MAxConn = ", config.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			//创建一个connertion的对象
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
			s.OnConnStart(dealConn)
		}
	}()

}
func (s *Server) Stop() {
	s.connManger.ClearConn()
}
func (s *Server) Server() {
	s.Start()
	select {}
}
//添加路由方法
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add msgID Success", msgId)
}
func (s *Server) GetConnMgr() ziface.IConnmanger {
	return s.connManger
}

//注册创建链接之后调用的hook函数的方法
func (s *Server) AddOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册销毁链接之前调用的hook函数的方法
func (s *Server) AddOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用创建链接之后调用的hook函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
if s.OnConnStart!=nil{
	fmt.Println("---> Call OnConnStart()...")
	s.OnConnStart(conn)
}
}

//调用函数销毁之前调用的hook函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop!=nil{
		fmt.Println("---> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
