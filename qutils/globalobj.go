package utils

import (
	"Q/qiface"
	"Q/qlog"
)

type GlobalObj struct {
	TcpServer qiface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	MaxMsgChanLen    uint32

	ConfigFilePath string

	LogDir        string
	LogFile       string
	LogDebugClose bool
}

var GlobalObject *GlobalObj

func init() {

	// default
	GlobalObject = &GlobalObj{
		Name:             "QServerApp",
		Version:          "V0.1",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		ConfigFilePath:   "./config/",
		LogDir:           ".",
		LogFile:          "log",
		LogDebugClose:    false,
	}

	if GlobalObject.LogFile != "" {
		qlog.SetLogFile(GlobalObject.LogDir, GlobalObject.LogFile)
	}
	if GlobalObject.LogDebugClose == true {
		qlog.CloseDebug()
	}
}
