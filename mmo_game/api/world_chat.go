package api

import (
	"Q/mmo_game/core"
	"Q/mmo_game/pb"
	"Q/qiface"
	"Q/qnet"
	"fmt"

	"github.com/golang/protobuf/proto"
)

//世界聊天 路由业务
type WorldChatApi struct {
	qnet.BaseRouter
}

func (*WorldChatApi) Handle(request qiface.IRequest) {
	//1. 将客户端传来的proto协议解码
	msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error ", err)
		return
	}

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//3. 根据pID得到player对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	//4. 让player对象发起聊天广播请求
	player.Talk(msg.Content)
}
