package grpc

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/config/domain"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/config"
	"log"
)

type Server struct {
	configService domain.ConfigService
	*kitGrpc.Server
	pb.UnimplementedConfigServiceServer
}

func New(configService domain.ConfigService) *Server {

	s := &Server{
		configService: configService,
	}

	// grpc server
	gs, err := kitGrpc.NewServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterConfigServiceServer(s.Srv, s)

	return s
}

func (s *Server) Get(ctx context.Context, rq *pb.ConfigRequest) (*pb.ConfigResponse, error) {

	if cfg, err := s.configService.Get(); err == nil {
		cfgj, err := json.Marshal(cfg)
		if err != nil {
			return nil, err
		}
		return &pb.ConfigResponse{Config: cfgj}, nil
	} else {
		return nil, err
	}

}

func (s *Server) ListenAsync() {

	go func() {
		grpc := s.configService.GrpcSettings()
		err := s.Server.Listen(grpc.Hosts[0], grpc.Port)
		if err != nil {
			log.Fatal(err)
		}
	}()

}
