package grpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	Srv     *grpc.Server
	Service string
}

func NewServer(service string) (*Server, error) {

	s := &Server{
		Service: service,
		Srv: grpc.NewServer(grpc_middleware.WithUnaryServerChain(
			ContextUnaryServerInterceptor())),
	}

	return s, nil
}

func (s *Server) Listen(host, port string) error {


	log.L().Cmp(s.Service).Pr("grpc").F(log.FF{"host": host, "port": port}).Inf("start listening")

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
