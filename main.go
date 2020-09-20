package main

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"log"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	log.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendBuffMsg(request.GetMsgID(), []byte("ping...ping...ping"))
	if err != nil {
		log.Println("error:", err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.Serve()
}
