package qiface

type IServer interface {
	Start() // start server

	Stop() // stop server

	Serve() // open server

	AddRouter(msgId uint32, router IRouter) // add handler for  msgID

	GetConnMgr() IConnManager // get conn manager

	SetOnConnStart(func(conn IConnection)) // set hook for conn start

	SetOnConnStop(func(conn IConnection)) // set hook for conn stop

	CallOnConnStart(conn IConnection) // call hook for conn start

	CallOnConnStop(conn IConnection) // call hook for conn stop

	DataPack() IDataPack
}
