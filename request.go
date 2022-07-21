package amethyst

import "github.com/phrara/amethyst/interval"

type Request struct {
	conn     interval.Connecter
	data     *interval.Packet
	protocol string
}

func NewRequest(conn interval.Connecter, data *interval.Packet, proto string) *Request {
	return &Request{
		conn:     conn,
		data:     data,
		protocol: proto,
	}
}

func (r *Request) Conn() interval.Connecter {
	return r.conn
}

func (r *Request) ConnID() uint {
	return r.conn.ConnID()
}

func (r *Request) Tag() uint32 {
	return r.data.Tag
}

func (r *Request) Len() uint32 {
	return r.data.Len
}

func (r *Request) Data() []byte {
	return r.data.Value
}

func (r *Request) Protocol() string {
	return r.protocol
}
