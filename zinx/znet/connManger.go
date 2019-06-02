package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)
//链接管理模块
type ConnManger struct {
	//管理所有链接集合
	connections map[uint32]ziface.IConnection
	//加读写锁
	connLock sync.RWMutex
}

//初始化方法
func NewConnManger() ziface.IConnmanger {
	return &ConnManger{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//添加链接
func (connMgr *ConnManger) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("Add connid = ", conn.GetConnID(), "to manager succ!!")
}

//删除链接
func (connMgr *ConnManger) Remove(connID uint32) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, connID)
	fmt.Println("Remove connid = ", connID, " from manager succ!!")
}

//根据链接ID得到链接
func (connMgr *ConnManger) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, OK := connMgr.connections[connID]; OK {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!")
	}
}

//得到目前服务器链接总个数
func (connMgr *ConnManger) Len() int {
	return len(connMgr.connections)
}

//清空全部链接的方法
func (connMgr *ConnManger) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//遍历删除
	for connID, conn := range connMgr.connections {
		//将全部的conn 关闭
		conn.Stop()

		//删除链接
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All Conections succ! conn num = ", connMgr.Len())
}
