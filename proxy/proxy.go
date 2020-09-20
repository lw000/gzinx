package proxy

import (
	"context"
	"github.com/aceld/zinx/znet"
	"io"
	"log"
	"net"
	"time"
)

type Service struct {
	host   string
	port   string
	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc
	closed bool
}

func CreateProxy(host, port string) *Service {
	return &Service{
		host: host,
		port: port,
	}
}

func (s *Service) Start() error {
	var err error
	s.conn, err = net.DialTimeout("tcp", net.JoinHostPort(s.host, s.port), time.Second*15)
	if err != nil {
		return err
	}
	s.closed = false
	go s.readLoop()
	return nil
}

func (s *Service) Stop() {
	_ = s.conn.Close()
}

func (s *Service) SendMessage(id uint32, data []byte) error {
	dp := znet.NewDataPack()
	msg, _ := dp.Pack(znet.NewMsgPackage(id, data))
	_, err := s.conn.Write(msg)
	if err != nil {
		log.Println("write error err ", err)
		return err
	}
	return nil
}

func (s *Service) readLoop() {
	defer func() {
		if p := recover(); p != nil {
			log.Println(p)
		}
	}()

	for {
		dp := znet.NewDataPack()

		// 先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(s.conn, headData) // ReadFull 会把msg填充满为止
		if err != nil {
			log.Println("read head error")
			break
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
			_, err := io.ReadFull(s.conn, msg.Data)
			if err != nil {
				log.Println("server unpack data err:", err)
				return
			}

			log.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}
	}
}

func (s *Service) keepLive() {
	if s.closed {

	}
}
