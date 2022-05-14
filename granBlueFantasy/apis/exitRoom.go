package apis

import (
	"Q/granBlueFantasy/core"
	accountMsg "Q/granBlueFantasy/pb"
	"Q/qiface"
	"Q/qnet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type ExitRoomApi struct {
	qnet.BaseRouter
}

func (*ExitRoomApi) Handle(request qiface.IRequest) {
	reqMsg := &accountMsg.LeaveRoomReq{}

	err := proto.Unmarshal(request.GetData(), reqMsg)

	if err != nil {
		fmt.Println("Leave Unmarshal error ", err)
		return
	}

	core.TheWorldManager.LeaveRoom(reqMsg.GetRoomId(), request.GetConnection().GetConnID())

	resMsg := &accountMsg.LeaveRoomRes{
		Res: &accountMsg.BaseResponse{
			Code:    0,
			Message: "离开房间成功",
		},
	}

	resMsgProto, err := proto.Marshal(resMsg)
	if err != nil {
		fmt.Println("LeaveRoom Marshal error ", err)
		return
	}

	err = request.GetConnection().SendMsg(10012, resMsgProto)
	if err != nil {
		return
	}
}
