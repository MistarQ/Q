package qiface

import (
	"context"
	"net"
)

type IConnection interface {
	Start() // start connection

	Stop() // stop connection

	Context() context.Context // return ctx

	GetTcpConnection() *net.TCPConn // get tcp conn

	GetConnID() uint32 // get connection ID

	GetRemoteAddr() net.Addr // get remote addr

	SendMsg(msgId uint32, data []byte) error // send msg directly

	SendBuffMsg(msgId uint32, data []byte) error // send msg with buffer

	SetProperty(key string, value any) // set conn property

	GetProperty(key string) (any, error) // get conn property

	RemoveProperty(key string) // remove conn property
}

type HandleFunc func(*net.TCPConn, []byte, int) error
