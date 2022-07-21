package amethyst

import (
	"fmt"
	"github.com/phrara/amethyst/config"
	"github.com/phrara/amethyst/interval"
	"net"
)

// Server Server implement

type Server struct {
	ServerID   string
	Protocol   string
	IP         string
	Port       string
	routerMngr *RouterManager
	connMngr   *ConnManager
	middleWare *middleWare
}

func generate(name, ip, port string, protocol string) *Server {

	return &Server{
		ServerID:   name,
		Protocol:   protocol,
		IP:         ip,
		Port:       port,
		routerMngr: NewRouterManager(),
		connMngr:   NewConnManager(),
		middleWare: &middleWare{
			ware: make(map[int][]func(connecter interval.Connecter), 2),
		},
	}
}

func (s *Server) Init() bool {
	fmt.Printf("\u001B[1;35m%s\u001B[0m\n", s.ServerID)
	c := make(chan error, 1)
	switch s.Protocol {
	case "tcp":
		go s.initTCPServer(c)
		if err := <-c; err != nil {
			return false
		}
	case "udp":
		go s.initUDPServer(c)
		if err := <-c; err != nil {
			return false
		}
	default:
		return false
	}
	return true
}

func (s *Server) Run() {
	if b := s.Init(); b {
		info := fmt.Sprintf("Amethyst server run at %s: %s (%s)", s.IP, s.Port, s.Protocol)
		fmt.Printf("\x1b[1;34m%s\x1b[0m\n", info)

		select {}
	} else {
		return
	}

}

func (s *Server) Close() {
	s.routerMngr.wp.Shut()
	s.connMngr.Clear()
}

func (s *Server) Route(tag uint32, r Router) {
	s.routerMngr.AppendRouter(tag, r)
}

func (s *Server) AddrString() string {
	return s.IP + ":" + s.Port
}

func (s *Server) initTCPServer(c chan<- error) {
	addr, err := net.ResolveTCPAddr(string(s.Protocol), s.AddrString())
	if err != nil {
		c <- err
		return
	}
	tcpListener, err := net.ListenTCP(string(s.Protocol), addr)
	if err != nil {
		c <- err
		return
	}
	c <- nil
	id := 0
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Get a client: ", conn.RemoteAddr().String())

		// check whether current conn num is lager than MaxConn
		if s.connMngr.Sum() >= config.Global.MaxConn {
			conn.Write([]byte("The server is busy. Please try again later"))
			conn.Close()
			continue
		}
		tcpConn := NewTCPConnection(conn, uint(id), s.routerMngr, s.middleWare)
		tcpConn.Start()
		s.connMngr.Append(tcpConn)
		id++
	}
}

func (s *Server) initUDPServer(c chan<- error) {
	c <- nil
}

// Register
// when is 0: middleWare will be called after conn established;
// when is 1: middleWare will be called before conn closed
func (s *Server) Register(when int, middleWare ...func(connecter interval.Connecter)) {
	if s.middleWare.ware[when] == nil {
		s.middleWare.ware[when] = make([]func(connecter interval.Connecter), 0)
	}
	for _, f := range middleWare {
		s.middleWare.ware[when] = append(s.middleWare.ware[when], f)
	}
}

type middleWare struct {
	ware map[int][]func(connecter interval.Connecter)
}
