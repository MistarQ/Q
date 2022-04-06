package qnet

import (
	"Q/qiface"
	"Q/qutils"
	"fmt"
	"net"
)

type Server struct {
	Name string

	IPVersion string

	IP string

	Port int

	MsgHandler qiface.IMsgHandler

	ConnManager qiface.IConnManager

	OnConnStart func(conn qiface.IConnection)

	OnConnStop func(conn qiface.IConnection)

	dataPack qiface.IDataPack
}

func NewServer(opts ...Option) qiface.IServer {

	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.Port,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
		dataPack:    NewDataPack(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	fmt.Printf("[Q] Serer Name: %s, listener at IP : %s, Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.Port)
	fmt.Printf("[Q] Version %s, MaxConn:%d, MaxPacketSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {

		s.MsgHandler.StartWorkerPool()

		// get tcp addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error : ", err)
			return
		}
		// listen server addr
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start Q server succeed, ", s.Name, "succeed, Listening...")

		// TODO cid应为steamID
		var cid uint32 = 0

		// accept client
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应错误
				fmt.Println("Too many connections, MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将服务器资源，状态，已经开辟的链接信息回收释放
	fmt.Println("[STOP] Q Server name: ", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	// TODO 扩展

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router qiface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Printf("Add Router Succeed, msgID=%d", msgID)
}

func (s *Server) GetConnMgr() qiface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(hookFunc func(conn qiface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn qiface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn qiface.IConnection) {
	if s.OnConnStart == nil {
		return
	}
	fmt.Println("->CallOnConnStart()...")
	s.OnConnStart(conn)
}

func (s *Server) CallOnConnStop(conn qiface.IConnection) {
	if s.OnConnStop == nil {
		return
	}
	fmt.Println("->CallOnConnStop()...")
	s.OnConnStop(conn)
}

func (s *Server) DataPack() qiface.IDataPack {
	return s.dataPack
}
