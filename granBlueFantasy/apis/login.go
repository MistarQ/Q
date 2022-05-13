package apis

import (
	accountMsg "Q/granBlueFantasy/pb"
	"Q/qiface"
	"Q/qnet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type LoginApi struct {
	qnet.BaseRouter
}

func (*LoginApi) Handle(request qiface.IRequest) {
	reqMsg := &accountMsg.LoginReq{}

	err := proto.Unmarshal(request.GetData(), reqMsg)

	if err != nil {
		fmt.Println("Login Unmarshal error ", err)
		return
	}

	// TODO 验证登陆信息

	resMsg := &accountMsg.LoginRes{
		Result: 1,
	}

	resMsgProto, err := proto.Marshal(resMsg)
	if err != nil {
		fmt.Println("Login Marshal error ", err)
		return
	}

	err = request.GetConnection().SendMsg(10002, resMsgProto)
	if err != nil {
		return
	}

}
