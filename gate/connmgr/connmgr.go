package connmgr

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type ConnItem struct {
	conn   *net.TCPConn
	connID uint32
	status int
	userId uint32
}

type ConnMgr struct {
	connLock sync.RWMutex
	conns    map[uint32]*ConnItem
}

func NewConnItem(conn *net.TCPConn, connID uint32) *ConnItem {
	return &ConnItem{
		conn:   conn,
		connID: connID,
	}
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		conns: make(map[uint32]*ConnItem),
	}
}

func (c *ConnItem) GetConnID() uint32 {
	return c.connID
}

func (c *ConnItem) GetConn() *net.TCPConn {
	return c.conn
}

func (c *ConnItem) Read(buf []byte) error {
	_, err := io.ReadFull(c.conn, buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConnItem) Write(data []byte) error {
	_, err := c.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConnItem) Close() {
	err := c.conn.Close()
	if err != nil {
	}
	c.connID = 0
}

func (c *ConnMgr) Add(conn *ConnItem) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.conns[conn.GetConnID()] = conn
}

func (c *ConnMgr) Remove(conn *ConnItem) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.conns, conn.GetConnID())

	log.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", c.Len())
}

func (c *ConnMgr) Clear() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.conns {
		_ = conn.GetConn().Close()
		delete(c.conns, connID)
	}
	log.Println("Clear All Connections successfully: conn num = ", c.Len())
}

// 利用ConnID获取链接
func (c *ConnMgr) Get(connID uint32) (*net.TCPConn, error) {
	// 保护共享资源Map 加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.conns[connID]; ok {
		return conn.GetConn(), nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnMgr) Len() int {
	return len(c.conns)
}
