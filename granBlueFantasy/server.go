package main

import (
	"Q/granBlueFantasy/apis"
	"Q/granBlueFantasy/db"
	"Q/qiface"
	"Q/qnet"
)

func OnConnectionAdd(conn qiface.IConnection) {

}

func OnConnectionLost(conn qiface.IConnection) {

}

func main() {

	s := qnet.NewServer()

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	s.AddRouter(10001, &apis.LoginApi{})
	s.AddRouter(10005, &apis.RoomListApi{})
	s.AddRouter(10007, &apis.CreateRoomApi{})
	s.AddRouter(10009, &apis.JoinRoomApi{})
	s.AddRouter(10015, &apis.RegisterApi{})
	s.AddRouter(20001, &apis.AttackApi{})

	db.InitDB()

	s.Serve()
}
