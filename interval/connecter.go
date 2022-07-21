package interval

import "net"

type Connecter interface {
	Start()
	Close()
	RemoteAddr() net.Addr
	Send(tag uint32, data []byte) error
	ConnID() uint
	State() bool
}
