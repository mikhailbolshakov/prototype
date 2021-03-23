package grpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

type Server struct {
	healthpb.HealthServer
	Srv     *grpc.Server
	Service string
	logger  log.CLoggerFunc
}

func NewServer(service string, logger log.CLoggerFunc) (*Server, error) {

	s := &Server{
		Service:      service,
		Srv:          grpc.NewServer(grpc_middleware.WithUnaryServerChain(ContextUnaryServerInterceptor()), grpc_middleware.WithStreamServerChain()),
		HealthServer: NewHealthServer(),
		logger:       logger,
	}

	healthpb.RegisterHealthServer(s.Srv, s)

	return s, nil
}

func (s *Server) Listen(host, port string) error {

	l := s.logger().Cmp(s.Service).Pr("grpc").Mth("listen").F(log.FF{"host": host, "port": port}).Inf("start listening")

	lis, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		l.E(err).Err("net.listen error")
		return err
	}

	err = s.Srv.Serve(lis)
	if err != nil {
		l.E(err).Err("serve error")
		return err
	}

	return nil

}

func (s *Server) Close() {
	s.Srv.Stop()
}
