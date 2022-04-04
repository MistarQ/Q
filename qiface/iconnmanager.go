package qiface

type IConnManager interface {
	Add(conn IConnection) // add conn

	Remove(conn IConnection) // remove conn

	Get(connID uint32) (IConnection, error) // get conn by connID

	Len() int // get conn nums

	ClearConn() // close all conn
}
