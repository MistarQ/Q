package qnet

import "Q/qiface"

type Request struct {
	conn qiface.IConnection

	msg qiface.IMessage
}

func (r *Request) GetConnection() qiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
