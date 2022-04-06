package qnet

import (
	"Q/qiface"
	"Q/qutils"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Connection struct {
	Server qiface.IServer

	Conn *net.TCPConn

	ConnId uint32

	isClosed bool

	msgBuffChan chan []byte

	MsgHandler qiface.IMsgHandler

	sync.RWMutex

	properties map[string]any

	propertiesLock sync.RWMutex

	ctx context.Context

	cancel context.CancelFunc
}

func NewConnection(server qiface.IServer, conn *net.TCPConn, connID uint32, msgHandler qiface.IMsgHandler) qiface.IConnection {
	c := &Connection{
		Server:      server,
		Conn:        conn,
		ConnId:      connID,
		isClosed:    false,
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		MsgHandler:  msgHandler,
		properties:  nil,
	}

	c.Server.GetConnMgr().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("[Reader is exit] connID = ", c.ConnId, ", remote addr = ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			headData := make([]byte, c.Server.DataPack().GetHeadLen())
			_, err := io.ReadFull(c.GetTcpConnection(), headData)
			if err != nil {
				fmt.Println("Read message head error", err)
				return
			}

			msg, err := c.Server.DataPack().Unpack(headData)
			if err != nil {
				fmt.Println("unpack error", err)
				return
			}

			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				_, err := io.ReadFull(c.GetTcpConnection(), data)
				if err != nil {
					fmt.Println("read msg data error", err)
					return
				}
			}
			msg.SetMsgData(data)

			req := Request{conn: c,
				msg: msg}

			if utils.GlobalObject.WorkerPoolSize > 0 {
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				go c.MsgHandler.DoMsgHandler(&req)
			}
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println("[conn Writer exit] ", c.GetRemoteAddr().String())

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if !ok {
				fmt.Println("msg buff chan is closed")
				// break
				return
			}
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error", err)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID = ", c.ConnId)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.StartReader()
	go c.StartWriter()

	c.Server.CallOnConnStart(c)
	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnId)
	c.cancel()
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	//将data封包，并且发送
	dp := c.Server.DataPack()
	msg, err := dp.Pack(NewMsgPack(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	_, err = c.Conn.Write(msg)
	return err
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	dp := c.Server.DataPack()

	binaryMsg, err := dp.Pack(NewMsgPack(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id=", msgID)
		return errors.New("pack error msg")
	}
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- binaryMsg:
		return nil
	}
}

func (c *Connection) SetProperty(key string, value any) {
	c.propertiesLock.Lock()
	defer c.propertiesLock.Unlock()
	if c.properties == nil {
		c.properties = make(map[string]any)
	}
	c.properties[key] = value
}

func (c *Connection) GetProperty(key string) (any, error) {
	c.propertiesLock.RLock()
	defer c.propertiesLock.RUnlock()

	if value, ok := c.properties[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertiesLock.Lock()
	defer c.propertiesLock.Unlock()

	delete(c.properties, key)
}

func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) finalizer() {
	c.Server.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()

	if c.isClosed {
		return
	}

	fmt.Println("Conn Stop()...ConnID = ", c.ConnId)

	_ = c.Conn.Close()

	c.Server.GetConnMgr().Remove(c)

	close(c.msgBuffChan)

	c.isClosed = true
}
