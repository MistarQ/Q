package apis

import (
	"Q/granBlueFantasy/core"
	accountMsg "Q/granBlueFantasy/pb"
	"Q/qiface"
	"Q/qnet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type CreateRoomApi struct {
	qnet.BaseRouter
}

func (*CreateRoomApi) Handle(request qiface.IRequest) {
	reqMsg := &accountMsg.CreateRoomReq{}

	err := proto.Unmarshal(request.GetData(), reqMsg)

	if err != nil {
		fmt.Println("CreateRoom Unmarshal error ", err)
		return
	}

	// TODO
	core.TheWorldManager.CreateRoom(reqMsg.BossId, 0)

	resMsg := &accountMsg.CreateRoomRes{
		Result: 1,
	}

	resMsgProto, err := proto.Marshal(resMsg)
	if err != nil {
		fmt.Println("CreateRoom Marshal error ", err)
		return
	}

	err = request.GetConnection().SendMsg(10008, resMsgProto)
	if err != nil {
		return
	}

}
