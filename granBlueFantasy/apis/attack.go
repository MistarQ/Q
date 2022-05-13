package apis

import (
	"Q/granBlueFantasy/core"
	accountMsg "Q/granBlueFantasy/pb"
	"Q/qiface"
	"Q/qnet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type AttackApi struct {
	qnet.BaseRouter
}

func (*AttackApi) Handle(request qiface.IRequest) {
	reqMsg := &accountMsg.AttackReq{}

	err := proto.Unmarshal(request.GetData(), reqMsg)

	if err != nil {
		fmt.Println("Attack Unmarshal error ", err)
		return
	}

	boss := core.TheWorldManager.GetRoom(reqMsg.GetRoomId()).GetBoss()
	boss.Hp -= 1
	resMsg := &accountMsg.AttackRes{
		RoomId:   reqMsg.GetRoomId(),
		PlayerId: reqMsg.GetPlayerId(),
		Boss:     boss,
	}

	resMsgProto, err := proto.Marshal(resMsg)
	if err != nil {
		fmt.Println("Attack Marshal error ", err)
		return
	}

	err = request.GetConnection().SendMsg(20002, resMsgProto)
	if err != nil {
		return
	}

}
