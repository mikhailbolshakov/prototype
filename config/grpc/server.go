package grpc

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/config/domain"
	"gitlab.medzdrav.ru/prototype/config/logger"
	"gitlab.medzdrav.ru/prototype/config/meta"
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
	gs, err := kitGrpc.NewServer(meta.ServiceCode, logger.LF())
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterConfigServiceServer(s.Srv, s)

	return s
}

func (s *Server) Get(ctx context.Context, rq *pb.ConfigRequest) (*pb.ConfigResponse, error) {

	if cfg, err := s.configService.Get(ctx); err == nil {
		cfgj, err := json.Marshal(cfg)
		if err != nil {
			return nil, err
		}
		return &pb.ConfigResponse{Config: cfgj}, nil
	} else {
		return nil, err
	}

}

func (s *Server) ListenAsync(ctx context.Context) {

	go func() {
		grpc := s.configService.GrpcSettings(ctx)
		err := s.Server.Listen(grpc.Host, grpc.Port)
		if err != nil {
			log.Fatal(err)
		}
	}()

}
