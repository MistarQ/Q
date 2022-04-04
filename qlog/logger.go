package qlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	LOG_MAX_BUF = 1024 * 1024
)

//日志头部信息标记位，采用bitmap方式，用户可以选择头部需要哪些标记位被打印
const (
	BitDate         = 1 << iota                            //日期标记位  2019/01/23
	BitTime                                                //时间标记位  01:23:12
	BitMicroSeconds                                        //微秒级标记位 01:23:12.111222
	BitLongFile                                            //完整文件名称
	BitShortFile                                           //最后文件名   server.go
	BitLevel                                               //当前日志级别： 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
	BitStdFlag      = BitDate | BitTime                    //标准头部日志格式
	BitDefault      = BitLevel | BitShortFile | BitStdFlag //默认日志头部格式
)

//日志级别
const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogError
	LogPanic
	LogFatal
)

//日志级别对应的显示字符串
var levels = []string{
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

type Logger struct {
	mu         sync.Mutex
	prefix     string
	flag       int
	out        io.Writer
	buf        bytes.Buffer
	file       *os.File
	debugClose bool
	callDepth  int
}

func NewLogger(out io.Writer, prefix string, flag int) *Logger {
	logger := &Logger{
		out:        out,
		prefix:     prefix,
		flag:       flag,
		file:       nil,
		debugClose: false,
		callDepth:  2,
	}
	runtime.SetFinalizer(logger, CleanLogger)
	return logger
}

func CleanLogger(logger *Logger) {
	logger.closeFile()
}

func (logger *Logger) closeFile() {
	if logger.file != nil {
		if err := logger.file.Close(); err != nil {
			fmt.Println("close logger err: ", err)
		}
		logger.file = nil
		logger.out = os.Stderr
	}
}

func (logger *Logger) formatHeader(t time.Time, file string, line int, level int) {
	var buf *bytes.Buffer = &logger.buf
	//如果当前前缀字符串不为空，那么需要先写前缀
	if logger.prefix != "" {
		buf.WriteByte('<')
		buf.WriteString(logger.prefix)
		buf.WriteByte('>')
	}

	//已经设置了时间相关的标识位,那么需要加时间信息在日志头部
	if logger.flag&(BitDate|BitTime|BitMicroSeconds) != 0 {
		//日期位被标记
		if logger.flag&BitDate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			buf.WriteByte('/') // "2019/"
			itoa(buf, int(month), 2)
			buf.WriteByte('/') // "2019/04/"
			itoa(buf, day, 2)
			buf.WriteByte(' ') // "2019/04/11 "
		}

		//时钟位被标记
		if logger.flag&(BitTime|BitMicroSeconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			buf.WriteByte(':') // "11:"
			itoa(buf, min, 2)
			buf.WriteByte(':') // "11:15:"
			itoa(buf, sec, 2)  // "11:15:33"
			//微秒被标记
			if logger.flag&BitMicroSeconds != 0 {
				buf.WriteByte('.')
				itoa(buf, t.Nanosecond()/1e3, 6) // "11:15:33.123123
			}
			buf.WriteByte(' ')
		}

		// 日志级别位被标记
		if logger.flag&BitLevel != 0 {
			buf.WriteString(levels[level])
		}

		//日志当前代码调用文件名名称位被标记
		if logger.flag&(BitShortFile|BitLongFile) != 0 {
			//短文件名称
			if logger.flag&BitShortFile != 0 {
				short := file
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						//找到最后一个'/'之后的文件名称  如:/home/go/src/q.go 得到 "q.go"
						short = file[i+1:]
						break
					}
				}
				file = short
			}
			buf.WriteString(file)
			buf.WriteByte(':')
			itoa(buf, line, -1) //行数
			buf.WriteString(": ")
		}
	}
}

func (logger *Logger) OutPut(level int, s string) error {

	now := time.Now() // 得到当前时间
	var file string   //当前调用日志接口的文件名称
	var line int      //当前代码行数
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if logger.flag&(BitShortFile|BitLongFile) != 0 {
		logger.mu.Unlock()
		var ok bool
		//得到当前调用者的文件名称和执行到的代码行数
		_, file, line, ok = runtime.Caller(logger.callDepth)
		if !ok {
			file = "unknown-file"
			line = 0
		}
		logger.mu.Lock()
	}

	//清零buf
	logger.buf.Reset()
	//写日志头
	logger.formatHeader(now, file, line, level)
	//写日志内容
	logger.buf.WriteString(s)
	//补充回车
	if len(s) > 0 && s[len(s)-1] != '\n' {
		logger.buf.WriteByte('\n')
	}

	//将填充好的buf 写到IO输出上
	_, err := logger.out.Write(logger.buf.Bytes())
	return err
}

// ====> Debug <====
func (logger *Logger) Debugf(format string, v ...interface{}) {
	if logger.debugClose == true {
		return
	}
	_ = logger.OutPut(LogDebug, fmt.Sprintf(format, v...))
}

func (logger *Logger) Debug(v ...interface{}) {
	if logger.debugClose == true {
		return
	}
	_ = logger.OutPut(LogDebug, fmt.Sprintln(v...))
}

// ====> Info <====
func (logger *Logger) Infof(format string, v ...interface{}) {
	_ = logger.OutPut(LogInfo, fmt.Sprintf(format, v...))
}

func (logger *Logger) Info(v ...interface{}) {
	_ = logger.OutPut(LogInfo, fmt.Sprintln(v...))
}

// ====> Warn <====
func (logger *Logger) Warnf(format string, v ...interface{}) {
	_ = logger.OutPut(LogWarn, fmt.Sprintf(format, v...))
}

func (logger Logger) Warn(v ...interface{}) {
	_ = logger.OutPut(LogWarn, fmt.Sprintln(v...))
}

// ====> Error <====
func (logger *Logger) Errorf(format string, v ...interface{}) {
	_ = logger.OutPut(LogError, fmt.Sprintf(format, v...))
}

func (logger *Logger) Error(v ...interface{}) {
	_ = logger.OutPut(LogError, fmt.Sprintln(v...))
}

// ====> Fatal 需要终止程序 <====
func (logger *Logger) Fatalf(format string, v ...interface{}) {
	_ = logger.OutPut(LogFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (logger *Logger) Fatal(v ...interface{}) {
	_ = logger.OutPut(LogFatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// ====> Panic  <====
func (logger *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	_ = logger.OutPut(LogPanic, s)
	panic(s)
}

func (logger *Logger) Panic(v ...interface{}) {
	s := fmt.Sprintln(v...)
	_ = logger.OutPut(LogPanic, s)
	panic(s)
}

// ====> Stack  <====
func (logger *Logger) Stack(v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, LOG_MAX_BUF)
	n := runtime.Stack(buf, true) //得到当前堆栈信息
	s += string(buf[:n])
	s += "\n"
	_ = logger.OutPut(LogError, s)
}

//获取当前日志bitmap标记
func (logger *Logger) Flags() int {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	return logger.flag
}

//重新设置日志Flags bitMap 标记位
func (logger *Logger) ResetFlags(flag int) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.flag = flag
}

//添加flag标记
func (logger *Logger) AddFlag(flag int) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.flag |= flag
}

//设置日志的 用户自定义前缀字符串
func (logger *Logger) SetPrefix(prefix string) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.prefix = prefix
}

//设置日志文件输出
func (logger *Logger) SetLogFile(fileDir string, fileName string) {
	var file *os.File

	//创建日志文件夹
	_ = mkdirLog(fileDir)

	fullPath := fileDir + "/" + fileName
	if logger.checkFileExist(fullPath) {
		//文件存在，打开
		file, _ = os.OpenFile(fullPath, os.O_APPEND|os.O_RDWR, 0644)
	} else {
		//文件不存在，创建
		file, _ = os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	}

	logger.mu.Lock()
	defer logger.mu.Unlock()

	//关闭之前绑定的文件
	logger.closeFile()
	logger.file = file
	logger.out = file
}

func (logger *Logger) CloseDebug() {
	logger.debugClose = true
}

func (logger *Logger) OpenDebug() {
	logger.debugClose = false
}

// ================== 以下是一些工具方法 ==========

//判断日志文件是否存在
func (logger *Logger) checkFileExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func mkdirLog(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0775); err != nil {
			if os.IsPermission(err) {
				e = err
			}
		}
	}
	return
}

//将一个整形转换成一个固定长度的字符串，字符串宽度应该是大于0的
//要确保buffer是有容量空间的
func itoa(buf *bytes.Buffer, i int, wID int) {
	var u uint = uint(i)
	if u == 0 && wID <= 1 {
		buf.WriteByte('0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wID > 0; u /= 10 {
		bp--
		wID--
		b[bp] = byte(u%10) + '0'
	}

	// avoID slicing b to avoID an allocation.
	for bp < len(b) {
		buf.WriteByte(b[bp])
		bp++
	}
}
