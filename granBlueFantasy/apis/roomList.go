package apis

import (
	"Q/granBlueFantasy/core"
	accountMsg "Q/granBlueFantasy/pb"
	"Q/qiface"
	"Q/qnet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type RoomListApi struct {
	qnet.BaseRouter
}

func (*RoomListApi) Handle(request qiface.IRequest) {
	reqMsg := &accountMsg.RoomListReq{}

	err := proto.Unmarshal(request.GetData(), reqMsg)

	if err != nil {
		fmt.Println("RoomList Unmarshal error ", err)
		return
	}

	roomList := make([]*accountMsg.Room, len(core.TheWorldManager.Rooms))
	for _, room := range core.TheWorldManager.Rooms {
		roomList = append(roomList, room)
	}

	resMsg := &accountMsg.RoomListRes{
		RoomList: roomList,
	}

	resMsgProto, err := proto.Marshal(resMsg)
	if err != nil {
		fmt.Println("RoomList Marshal error ", err)
		return
	}

	err = request.GetConnection().SendMsg(10006, resMsgProto)
	if err != nil {
		return
	}

}
