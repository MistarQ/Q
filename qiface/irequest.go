package qiface

type IRequest interface {
	GetConnection() IConnection // get conn

	GetData() []byte // get data

	GetMsgID() uint32 // get msgID
}
