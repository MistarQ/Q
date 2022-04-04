package qnet

import (
	"Q/qiface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connections map[uint32]qiface.IConnection

	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]qiface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn qiface.IConnection) {
	connMgr.connLock.Lock()
	connMgr.connections[conn.GetConnID()] = conn
	connMgr.connLock.Unlock()
	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager succeed: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn qiface.IConnection) {
	connMgr.connLock.Lock()
	delete(connMgr.connections, conn.GetConnID())
	connMgr.connLock.Unlock()
	fmt.Println("connID = ", conn.GetConnID(), " removed from ConnManager succeed: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (qiface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND")
	}
}

func (connMgr *ConnManager) Len() int {
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())

}
