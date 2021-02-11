package grpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	Srv *grpc.Server
}

func NewServer() (*Server, error) {

	s := &Server{
		Srv: grpc.NewServer(grpc_middleware.WithUnaryServerChain(
			ContextUnaryServerInterceptor())),
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

func (s *Server) Close() {
	s.Srv.Stop()
}