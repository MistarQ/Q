package apis

import (
	"Q/granBlueFantasy/db"
	accountMsg "Q/granBlueFantasy/pb"
	"Q/qiface"
	"Q/qnet"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gorm.io/gorm"
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
	u := &db.User{
		Account:  reqMsg.GetAccount().GetAccount(),
		Password: reqMsg.GetAccount().GetPassword(),
	}

	resMsg := &accountMsg.LoginRes{
		Res: &accountMsg.BaseResponse{
			Code:    0,
			Message: "登陆成功",
		},
	}

	if err := db.Mysql.Take(u).Error; err != nil {
		resMsg.GetRes().Code = 1
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resMsg.GetRes().Message = "账号或密码错误"
		} else {
			resMsg.GetRes().Message = "登陆失败, " + err.Error()
		}
	} else {
		resMsg.Id = uint32(u.ID)
		request.GetConnection().SetConnID(uint32(u.ID))
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
