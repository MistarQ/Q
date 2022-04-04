package qiface

type IMessage interface {
	GetMsgId() uint32 // get msg id

	GetDataLen() uint32 // get len of data

	GetData() []byte // get data

	SetMsgId(uint322 uint32) // set msg id

	SetDataLen(uint322 uint32) // set len of data

	SetMsgData([]byte) // set data
}
