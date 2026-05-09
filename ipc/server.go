package ipc

import (
	"net"
)

type State struct {
	Query string
}

type Server struct {
	conn       *net.UDPConn
	StateChan  chan State
	ListenPort int
}

func NewServer() (*Server, error) {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0, // OS assigns random port
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	s := &Server{
		conn:       conn,
		StateChan:  make(chan State, 100),
		ListenPort: localAddr.Port,
	}

	return s, nil
}

func (s *Server) Start() {
	buf := make([]byte, 4096)
	for {
		n, _, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			break
		}

		query := string(buf[:n])
		// send non-blocking so we don't choke the UDP read loop
		select {
		case s.StateChan <- State{Query: query}:
		default:
		}
	}
}

func (s *Server) Close() {
	_ = s.conn.Close()
	close(s.StateChan)
}
