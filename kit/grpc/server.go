package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	Srv *grpc.Server
}

func NewGrpcServer() (*Server, error) {

	var opts []grpc.ServerOption

	s := &Server{
		Srv: grpc.NewServer(opts...),
	}

	return s, nil
}

func (s *Server) Listen(host, port string) error {

	log.Printf("GRPC server listening on %s:%s", host, port)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return err
	}

	err = s.Srv.Serve(lis)
	if err != nil {
		return err
	}

	return nil

}