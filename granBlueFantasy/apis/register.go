package apis

import (
	accountMsg "Q/granBlueFantasy/pb"
	"Q/qiface"
	"Q/qnet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type RegisterApi struct {
	qnet.BaseRouter
}

func (*RegisterApi) Handle(request qiface.IRequest) {
	reqMsg := &accountMsg.RegisterReq{}

	err := proto.Unmarshal(request.GetData(), reqMsg)

	if err != nil {
		fmt.Println("Register Unmarshal error ", err)
		return
	}

	resMsg := &accountMsg.RegisterRes{
		Result: 1,
	}

	resMsgProto, err := proto.Marshal(resMsg)
	if err != nil {
		fmt.Println("Register Marshal error ", err)
		return
	}

	err = request.GetConnection().SendMsg(10016, resMsgProto)
	if err != nil {
		return
	}

}
