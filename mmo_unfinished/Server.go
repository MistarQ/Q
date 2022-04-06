package main

import (
	"Q/qiface"
	"Q/qnet"
	"fmt"
)

type PingRouter struct {
	qnet.BaseRouter
}

func (this *PingRouter) Handle(request qiface.IRequest) {
	fmt.Println("Call Ping Router Handle...")
	fmt.Println("receive from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	qnet.BaseRouter
}

func (this *HelloRouter) Handle(request qiface.IRequest) {
	fmt.Println("Call Hello Router Handle...")
	fmt.Println("receive from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn qiface.IConnection) {
	fmt.Println("->DoConnectionBegin is called...")
	err := conn.SendMsg(202, []byte("DoConnection BEGIN"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Set conn Name, Home...")
	conn.SetProperty("Name", "QS")
	conn.SetProperty("GitHub", "github.com/MistarQ")
}

func DoConnectionEnd(conn qiface.IConnection) {
	fmt.Println("->DoConnectionEnd is called...")
	fmt.Println("conn ID = ", conn.GetConnID(), " is lost...")

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if home, err := conn.GetProperty("GitHub"); err == nil {
		fmt.Println("GitHub = ", home)
	}
}

func main() {

	s := qnet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionEnd)

	s.AddRouter(0, &PingRouter{})

	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
