package agent

import (
	"net"
	"net/rpc"
)

type Server struct {
	*rpc.Server
}

func (s *Server) setupRPC(addr string) error {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	_ = listener
}
