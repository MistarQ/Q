package apis

import (
	"Q/granBlueFantasy/core"
	accountMsg "Q/granBlueFantasy/pb"
	"Q/qiface"
	"Q/qnet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type JoinRoomApi struct {
	qnet.BaseRouter
}

func (*JoinRoomApi) Handle(request qiface.IRequest) {
	reqMsg := &accountMsg.EnterRoomReq{}

	err := proto.Unmarshal(request.GetData(), reqMsg)

	if err != nil {
		fmt.Println("JoinRoom Unmarshal error ", err)
		return
	}

	// TODO
	room := core.TheWorldManager.JoinRoom(reqMsg.RoomId, 0)

	resMsg := &accountMsg.EnterRoomRes{
		Code: 1,
		Room: room,
	}

	resMsgProto, err := proto.Marshal(resMsg)
	if err != nil {
		fmt.Println("JoinRoom Marshal error ", err)
		return
	}

	err = request.GetConnection().SendMsg(10010, resMsgProto)
	if err != nil {
		return
	}

}
