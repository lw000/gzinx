package main

import (
	"flag"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	// using standard library "flag" package
	flag.Int("flag_name", 1234, "help message for flag_name")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Println(err)
	}

	value := viper.GetInt("flag_name")
	log.Println(value)

	s := znet.NewServer()
	s.AddRouter(1, &PingRouter{})
	s.Serve()
}
