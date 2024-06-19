package grpcserver

import (
	"google.golang.org/grpc"
	"net"
)

const (
	defaultAddr = ":44044"
)

type Server struct {
	server *grpc.Server
	notify chan error
}

func NewServer(g *grpc.Server) (*Server, error) {
	tcpServer, err := net.Listen("tcp", defaultAddr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		server: g,
		notify: make(chan error, 1),
	}

	s.start(tcpServer)
	return s, nil
}

func (s *Server) start(lis net.Listener) {
	go func() {
		s.notify <- s.server.Serve(lis)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() {
	s.server.GracefulStop()
}
