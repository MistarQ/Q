package qnet

import "Q/qiface"

type Option func(s *Server)

func withPacket(dataPack qiface.IDataPack) Option {
	return func(s *Server) {
		s.dataPack = dataPack
	}
}
