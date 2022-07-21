package amethyst

import (
	"errors"
	"fmt"
	"github.com/phrara/amethyst/config"
	"github.com/phrara/amethyst/interval"
	"io"
	"net"
	"sync"
)

const (
	CLOSED  = false
	OPENING = true
)

type TCPConnection struct {
	// connection
	Conn *net.TCPConn
	// conn ID
	ID uint
	// conn stat
	Stat bool
	// Exit Flag
	ExitC chan bool
	// router mngr
	routerMngr *RouterManager
	// Message chan
	msgC chan []byte
	// mutex
	m sync.Mutex
	// middleWare
	middleWare *middleWare
}

func NewTCPConnection(conn *net.TCPConn, id uint, router *RouterManager, ware *middleWare) *TCPConnection {
	return &TCPConnection{
		Conn:       conn,
		ID:         id,
		routerMngr: router,
		Stat:       OPENING,
		ExitC:      make(chan bool, 1),
		msgC:       make(chan []byte),
		middleWare: ware,
	}
}

func (c *TCPConnection) Start() {
	fmt.Printf("TCPConn %d is opening\n", c.ID)
	go c.readData()
	go c.writeData()

	for _, f := range c.middleWare.ware[0] {
		f(c)
	}

}

func (c *TCPConnection) writeData() {
	defer c.Close()
	for {
		select {
		case data := <-c.msgC:
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		case <-c.ExitC:
			return
		}
	}
}

func (c *TCPConnection) readData() {
	defer c.Close()
	for {
		p := &interval.Packet{}
		header := make([]byte, interval.HEADER)
		if _, err := io.ReadFull(c.Conn, header); err != nil {
			fmt.Println("read header from client failed: ", err)
			break
		}
		if err := p.ParseHeader(header); err != nil {
			fmt.Println("parse header failed: ", err)
			break
		}

		if p.Len > 0 {
			val := make([]byte, p.Len)
			if _, err := io.ReadFull(c.Conn, val); err != nil {
				fmt.Println("read value from client failed: ", err)
				break
			}
			p.Value = val
		}
		fmt.Printf("recv a packet: %#v from %s\n", *p, c.RemoteAddr().String())
		req := NewRequest(c, p, "tcp")
		if config.Global.MaxPoolSize > 0 {
			err := c.routerMngr.HandleRequest(req)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			go c.routerMngr.Handle(req)
		}

	}
}

func (c *TCPConnection) Close() {
	c.m.Lock()
	if c.Stat == CLOSED {
		return
	} else {

		for _, f := range c.middleWare.ware[1] {
			f(c)
		}

		c.Stat = CLOSED
		c.Conn.Close()
		close(c.ExitC)
		close(c.msgC)
	}
	c.m.Unlock()
	fmt.Printf("TCPConn %d closed\n", c.ID)
}

func (c *TCPConnection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *TCPConnection) Send(tag uint32, data []byte) error {
	if c.Stat == CLOSED {
		return errors.New("conn is closed")
	}
	p := &interval.Packet{
		Tag:   tag,
		Len:   uint32(len(data)),
		Value: data,
	}
	wrap, err := p.Wrap()
	if err != nil {
		return err
	}
	c.msgC <- wrap
	return nil
}

func (c *TCPConnection) ConnID() uint {
	return c.ID
}

func (c *TCPConnection) State() bool {
	return c.Stat
}
