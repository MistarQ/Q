
Q是一个基于Golang的轻量级并发服务器框架



> **说明**:作用领域包括:后端模块的消息中转、长链接游戏服务器、Web框架中的消息处理插件等。定位是代码简洁，让更多的开发者迅速的了解框架的内脏细节并且可以快速基于Q开发一款适合自己企业场景的模块。



## 三、初探架构
框架来源 https://github.com/aceld/zinx

![1-Zinx框架.png](https://camo.githubusercontent.com/903d1431358fa6f4634ebaae3b49a28d97e23d77/68747470733a2f2f75706c6f61642d696d616765732e6a69616e7368752e696f2f75706c6f61645f696d616765732f31313039333230352d633735666636383232333362323533362e706e673f696d6167654d6f6772322f6175746f2d6f7269656e742f7374726970253743696d61676556696577322f322f772f31323430)


## 三、开发API文档

### 快速开始

#### server
主函数步骤比较精简，最多只需要3步即可。
1. 创建server句柄
2. 配置自定义路由及业务
3. 启动服务

```go
func main() {
	//1 创建一个server句柄
	s := qnet.NewServer()

	//2 配置路由
	s.AddRouter(0, &PingRouter{})

	//3 开启服务
	s.Serve()
}
```

其中自定义路由及业务配置方式如下：
```go
import (
    "Q/qiface"
    "Q/qnet"
    "fmt"
)

//ping test 自定义路由
type PingRouter struct {
	qnet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handler(request qiface.IRequest) {
	//先读取客户端的数据
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

    //再回写
	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}
```

#### client
消息处理采用，`[MsgLength]|[MsgID]|[Data]`的封包格式
```go
package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"Q/qnet"
)

/*
	模拟客户端
 */
func main() {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn,err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for n := 3; n >= 0; n-- {
		//发封包message消息
		dp := qnet.NewDataPack()
		msg, _ := dp.Pack(q.NewMsgPackage(0,[]byte("Client Test Message")))
		_, err := conn.Write(msg)
		if err !=nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*qnet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1*time.Second)
	}
}
```

### 配置文件 json 后续可修改为yaml并配置viper
```json
{
  "Name":"zinx v-0.10 demoApp",
  "Host":"127.0.0.1",
  "Port":7777,
  "MaxConn":3,
  "WorkerPoolSize":10,
  "LogDir": "./mylog",
  "LogFile":"q.log"
}
```

`Name`:服务器应用名称

`Host`:服务器IP

`Port`:服务器监听端口

`MaxConn`:允许的客户端链接最大数量

`WorkerPoolSize`:工作任务池最大工作Goroutine数量

`LogDir`: 日志文件夹

`LogFile`: 日志文件名称(如果不提供，则日志信息打印到Stderr)


### I.服务器模块Server
```go
  func NewServer () qiface.IServer 
```
创建一个服务器句柄，该句柄作为当前服务器应用程序的主枢纽，包括如下功能：

#### 1)开启服务
```go
  func (s *Server) Start()
```
#### 2)停止服务
```go
  func (s *Server) Stop()
```
#### 3)运行服务
```go
  func (s *Server) Serve()
```
#### 4)注册路由
```go
func (s *Server) AddRouter (msgId uint32, router qiface.IRouter) 
```
#### 5)注册链接创建Hook函数
```go
func (s *Server) SetOnConnStart(hookFunc func (qiface.IConnection))
```
#### 6)注册链接销毁Hook函数
```go
func (s *Server) SetOnConnStop(hookFunc func (qiface.IConnection))
```
### II.路由模块

```go
//实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {}

//这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化
func (br *BaseRouter)PreHandle(req ziface.IRequest){}
func (br *BaseRouter)Handle(req ziface.IRequest){}
func (br *BaseRouter)PostHandle(req ziface.IRequest){}
```


### III.链接模块
#### 1)获取原始的socket TCPConn
```go
  func (c *Connection) GetTCPConnection() *net.TCPConn 
```
#### 2)获取链接ID
```go
  func (c *Connection) GetConnID() uint32 
```
#### 3)获取远程客户端地址信息
```go
  func (c *Connection) RemoteAddr() net.Addr 
```
#### 4)发送消息
```go
  func (c *Connection) SendMsg(msgId uint32, data []byte) error 
  func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error
```
#### 5)链接属性
```go
//设置链接属性
func (c *Connection) SetProperty(key string, value interface{})

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error)

//移除链接属性
func (c *Connection) RemoveProperty(key string) 
```

---

