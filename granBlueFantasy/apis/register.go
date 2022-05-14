package apis

import (
	"Q/granBlueFantasy/db"
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

	u := &db.User{
		Account:  reqMsg.GetAccount().GetAccount(),
		Password: reqMsg.GetAccount().GetPassword(),
	}

	resMsg := &accountMsg.RegisterRes{
		Res: &accountMsg.BaseResponse{
			Code:    0,
			Message: "注册成功",
		},
	}

	if err := db.Mysql.Create(u).Error; err != nil {
		resMsg.GetRes().Code = 1
		resMsg.GetRes().Message = "注册失败, " + err.Error()
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
