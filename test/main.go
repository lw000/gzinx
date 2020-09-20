package main

import (
	"demo/gzinx/consts"
	TLogin "demo/gzinx/protos/msg"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"time"
)

func runClient() {
	log.Println("Client Test ... start")
	conn, err := net.Dial("tcp", "127.0.0.1:7800") // 127.0.0.1:7777
	if err != nil {
		log.Println("client start err, exit!")
		return
	}

	for n := 10; n >= 0; n-- {
		var req TLogin.ReqRegister
		req.Account = "1313123123"
		req.Password = "12311212312"
		data, err := proto.Marshal(&req)
		if err != nil {
			log.Println("write error err ", err)
			return
		}

		// 发封包message消息
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(consts.LoginServeCommandRegister, data))
		_, err = conn.Write(msg)
		if err != nil {
			log.Println("write error err ", err)
			return
		}
		//
		// // 先读出流中的head部分
		// headData := make([]byte, dp.GetHeadLen())
		// _, err = io.ReadFull(conn, headData) // ReadFull 会把msg填充满为止
		// if err != nil {
		// 	log.Println("read head error")
		// 	break
		// }
		// // 将headData字节流 拆包到msg中
		// msgHead, err := dp.Unpack(headData)
		// if err != nil {
		// 	log.Println("server unpack err:", err)
		// 	return
		// }
		//
		// if msgHead.GetDataLen() > 0 {
		// 	// msg 是有data数据的，需要再次读取data数据
		// 	msg := msgHead.(*znet.Message)
		// 	msg.Data = make([]byte, msg.GetDataLen())
		//
		// 	// 根据dataLen从io中读取字节流
		// 	_, err := io.ReadFull(conn, msg.Data)
		// 	if err != nil {
		// 		log.Println("server unpack data err:", err)
		// 		return
		// 	}
		//
		// 	switch msg.GetMsgId() {
		// 	case consts.LoginServeCommandRegister:
		// 		var ack TLogin.AckRegister
		// 		err = proto.Unmarshal(msg.GetData(), &ack)
		// 		if err != nil {
		// 			log.Println("proto unpack data err:", err)
		// 			return
		// 		}
		// 		log.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", ack.Code, ack.Msg, ack.Token)
		// 	case consts.LoginServeCommandLogin:
		// 	case consts.CommandError:
		// 	}
		// }

		time.Sleep(1 * time.Second)
	}
}

func main() {
	runClient()
}
