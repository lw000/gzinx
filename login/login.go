package main

import (
	"demo/gzinx/consts"
	TLogin "demo/gzinx/protos/msg"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"log"
)

type LoginRouter struct {
	znet.BaseRouter
}

func (p *LoginRouter) Handle(request ziface.IRequest) {
	// log.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	// err := request.GetConnection().SendBuffMsg(request.GetMsgID(), []byte("login...login...login"))
	// if err != nil {
	// 	log.Println("error:", err)
	// }

	var req TLogin.ReqRegister
	err := proto.Unmarshal(request.GetData(), &req)
	if err != nil {
		log.Println("proto unmarshal err ", err)
		return
	}
	log.Println("recv from client : ", req)

	var ack TLogin.AckRegister
	ack.Code = 1000
	ack.Msg = "login"
	ack.Token = "123123123213213121231231"
	data, err := proto.Marshal(&ack)
	if err != nil {
		log.Println("proto marshal err ", err)
		return
	}
	err = request.GetConnection().SendBuffMsg(request.GetMsgID(), data)
	if err != nil {
		log.Println("error:", err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(consts.LoginServeCommandRegister, &LoginRouter{})
	s.Serve()
}
