package qlog

import (
	"os"
)

var StdLogger = NewLogger(os.Stderr, "", BitDefault)

//Flags 获取StdQLog 标记位
func Flags() int {
	return StdLogger.Flags()
}

//ResetFlags 设置StdQLog标记位
func ResetFlags(flag int) {
	StdLogger.ResetFlags(flag)
}

//AddFlag 添加flag标记
func AddFlag(flag int) {
	StdLogger.AddFlag(flag)
}

//SetPrefix 设置StdQLog 日志头前缀
func SetPrefix(prefix string) {
	StdLogger.SetPrefix(prefix)
}

//SetLogFile 设置StdQLog绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	StdLogger.SetLogFile(fileDir, fileName)
}

//CloseDebug 设置关闭debug
func CloseDebug() {
	StdLogger.CloseDebug()
}

//OpenDebug 设置打开debug
func OpenDebug() {
	StdLogger.OpenDebug()
}

//Debugf ====> Debug <====
func Debugf(format string, v ...interface{}) {
	StdLogger.Debugf(format, v...)
}

//Debug Debug
func Debug(v ...interface{}) {
	StdLogger.Debug(v...)
}

//Infof ====> Info <====
func Infof(format string, v ...interface{}) {
	StdLogger.Infof(format, v...)
}

//Info -
func Info(v ...interface{}) {
	StdLogger.Info(v...)
}

// ====> Warn <====
func Warnf(format string, v ...interface{}) {
	StdLogger.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdLogger.Warn(v...)
}

// ====> Error <====
func Errorf(format string, v ...interface{}) {
	StdLogger.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdLogger.Error(v...)
}

// ====> Fatal 需要终止程序 <====
func Fatalf(format string, v ...interface{}) {
	StdLogger.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdLogger.Fatal(v...)
}

// ====> Panic  <====
func Panicf(format string, v ...interface{}) {
	StdLogger.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdLogger.Panic(v...)
}

// ====> Stack  <====
func Stack(v ...interface{}) {
	StdLogger.Stack(v...)
}

func init() {
	//因为StdQLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	//一般的QLogger对象 calldDepth=2, StdQLog的calldDepth=3
	StdLogger.callDepth = 3
}
