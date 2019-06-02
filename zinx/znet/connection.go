package znet

//链接模块
import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/config"
	"zinx/ziface"
)

//具体的TCP链接模块
type Connection struct {
	//当前链接属于哪一个server创建
	server ziface.IServer
	//原生的socket套接字
	Conn *net.TCPConn
	//链接ID
	ConnID uint32
	//链接是否关闭
	isClosed bool
	//handleAPI ziface.HandleFunc
	//Router ziface.IRouter
	MsgHandler ziface.IMsgHandler
	//添加一个reader和writer通信的chanel
	msgChan chan []byte
	//添加chan用来reader通知writer conn已经关闭 需要推出的消息
	writerExitChan chan bool
	//当前链接模块所具备的一些属性集合
	property  map[string]interface{}

	//保护当前property的锁
	propertyLock sync.RWMutex
}

//初始化链接方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) ziface.IConnection {
	c := &Connection{
		server: server,
		Conn:   conn,
		ConnID: connID,
		//handleAPI: callback_api,
		//Router:   router,
		MsgHandler:     msgHandler,
		isClosed:       false,
		msgChan:        make(chan []byte),
		writerExitChan: make(chan bool),
		property:make(map[string]interface{}),
	}
	c.server.GetConnMgr().Add(c)
	return c

}

//建立针对链接读业务的方法
func (c *Connection) StartRead() {
	fmt.Println("read go is startin...")
	defer fmt.Println("connID..", c.ConnID, "Reader is exit, remote addr is = ", c.GetRemoteAddr().String())
	defer c.Stop()
	for {
		/*buf := make([]byte,config.GlobalObject.MaxPackageSize)
		cnt, err := c.Conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read error..", err)
			continue
		} else if cnt == 0 {
			fmt.Println("read stop...")
			break
		}*/
		//创建拆包封包的对象
		dp := NewDataPack()
		//读取客户端消息的头部
		headdata := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.Conn, headdata)
		if err != nil {
			fmt.Println("Readfull headlen err", err)
			break
		}
		//根据头部 获取数据的长度 进行第二次读取
		msg, err := dp.UnPack(headdata)
		if err != nil {
			fmt.Println("unpack err", err)
			break
		}
		//根据长度再次读取
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.Conn, data)
			if err != nil {
				fmt.Println("readfull msglen err", err)
				break
			}
		}
		msg.SetData(data)
		//将数据和链接进行绑定
		req := NewRequest(c, msg)
		//go c.MsgHandler.DoMsgHandler(req)
		//将req交给worker工作池来处理
		if config.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
		/*err = c.handleAPI(req)
		if err != nil {
			fmt.Println("connID", c.ConnID, "handele is err", err)
			break
		}*/
	}
}

//创建writer的goroutine 专门给客户端发送消息
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine start...")
	defer fmt.Println("[Writer Goroutine Stop...]")
	for {
		select { //通过channel进行数据传输
		case data := <-c.msgChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("send data err", err)
				return
			}
			//reader关闭 writer关闭
		case <-c.writerExitChan:
			return

		}
	}
}

//启动连接
func (c *Connection) Start() {
	fmt.Println("Conn Start（）  ... id = ", c.ConnID)
	//开启读操作
	go c.StartRead()
	//开启写操作
	go c.StartWriter()
}

//停止链接
func (c *Connection) Stop() {
	c.server.CallOnConnStop(c)

	fmt.Println("connId stop=", c.ConnID)
	//回收工作
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//通知writer链接已经关闭
	c.writerExitChan <- true
	_ = c.Conn.Close()

	//将当前链接从管理链接模块中删除
	c.server.GetConnMgr().Remove(c.ConnID)
	//释放资源
	close(c.msgChan)
	close(c.writerExitChan)
}

//获取链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取conn原生套接字
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取远程客户IP端地址
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据给客户端
func (c *Connection) Send(msgId uint32, msgData []byte) error {
	//判断用户在线状态
	if c.isClosed == true {
		return errors.New("Connection closed ..send Msg ")
	}
	//封装 初始化封包 拆包方法
	dp := NewDataPack()
	//封包 把数据长度 msgID 数据内容写入
	binarymsg, err := dp.Pack(NewMessage(msgId, msgData))
	if err != nil {
		fmt.Println("pack err", err)
		return err
	}
	/*_, err = c.Conn.Write(binarymsg)
	if err != nil {
		fmt.Println("send buf error")
		return err
	}*/
	//通过channel传输数据
	c.msgChan <- binarymsg
	return nil
}
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//添加一个链接属性
	c.property[key] = value
}

//获取属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	//读取属性
	if value, OK := c.property[key]; OK {
		return value, nil
	} else {
		return nil, errors.New("no property found" + key)
	}
}

//删除属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//删除属性
	delete(c.property, key)
}
