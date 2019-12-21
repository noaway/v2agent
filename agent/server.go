package agent

import (
	"net/rpc"
)

type Server struct {
	*rpc.Server
}

func (s *Server) setupRPC(addr string) error {
	return nil
}
