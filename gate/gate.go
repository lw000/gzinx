package main

import (
	"demo/gzinx/gate/connmgr"
	"demo/gzinx/proxy"
	"errors"
	"fmt"
	"github.com/aceld/zinx/znet"
	"log"
	"net"
)

var (
	conns      *connmgr.ConnMgr
	loginProxy *proxy.Service
)

func forwardMessage(msg *znet.Message) error {
	if msg.GetMsgId() < 1000 {
		// 网关服务消息
	} else if 1000 < msg.GetMsgId() && msg.GetMsgId() < 2000 {
		// 大厅服务消息
	} else if 2000 < msg.GetMsgId() && msg.GetMsgId() < 3000 {
		// 游戏服务消息
	} else if 3000 < msg.GetMsgId() && msg.GetMsgId() < 4000 {
		// 登录服务消息
		err := loginProxy.SendMessage(msg.GetMsgId(), msg.GetData())
		if err != nil {
			log.Println(err)
		}
	} else {
		return errors.New("unknown msgId")
	}
	return nil
}

func dealClientConn(connItem *connmgr.ConnItem) {
	defer func() {
		conns.Remove(connItem)
		connItem.Close()
	}()

	for {
		// 创建数据包
		dp := znet.NewDataPack()

		// 先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		err := connItem.Read(headData) // ReadFull 会把msg填充满为止
		if err != nil {
			log.Println("read head error")
			return
		}

		// 将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			log.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// 根据dataLen从io中读取字节流
			err := connItem.Read(msg.Data)
			if err != nil {
				log.Println("read data err:", err)
				return
			}

			err = forwardMessage(msg)
			if err != nil {
				_ = connItem.Write([]byte("unknown msgID"))
				return
			}
		}
	}
}

func runGateService(port string) {
	addr, err := net.ResolveTCPAddr("tcp4", net.JoinHostPort("0.0.0.0", port))
	if err != nil {
		log.Panicln(err)
		return
	}

	var listener *net.TCPListener
	listener, err = net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Panicln(err)
		return
	}

	// 已经监听成功
	log.Println("start server success, now listening...")

	var connId uint32 = 1
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err ", err)
			continue
		}
		log.Println("Get conn remote addr = ", conn.RemoteAddr().String())

		// 超过最大连接，那么则关闭此新的连接
		if conns.Len() >= 5000 {
			_ = conn.Close()
			continue
		}

		connItem := connmgr.NewConnItem(conn, connId)
		conns.Add(connItem)
		connId++
		go dealClientConn(connItem)
	}
}

func main() {
	conns = connmgr.NewConnMgr()
	loginProxy = proxy.CreateProxy("127.0.0.1", "7900")
	err := loginProxy.Start()
	if err != nil {
		log.Panicln(err)
	}
	runGateService("7800")
}
