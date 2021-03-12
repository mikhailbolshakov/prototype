package config

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/config/domain"
	"gitlab.medzdrav.ru/prototype/config/grpc"
	"gitlab.medzdrav.ru/prototype/config/logger"
	"gitlab.medzdrav.ru/prototype/config/meta"
	"gitlab.medzdrav.ru/prototype/kit/service"
)

type serviceImpl struct {
	service.Cluster
	domain domain.ConfigService
	grpc *grpc.Server
}

func New() service.Service {

	s := &serviceImpl{
		Cluster: service.NewCluster(logger.LF(), meta.Meta),
	}
	s.domain = domain.New()
	s.grpc = grpc.New(s.domain)
	return s
}

func (s *serviceImpl) GetCode() string {
	return meta.Meta.ServiceCode()
}

func (s *serviceImpl) Init(ctx context.Context) error {

	// load config
	if err := s.domain.Load(ctx); err != nil {
		return err
	}

	// get config
	cfg, err := s.domain.Get(ctx)
	if err != nil {
		return err
	}

	// set logging params
	srvCfg, ok := cfg.Services["cfg"]
	if ok {
		logger.Logger.SetLevel(srvCfg.Log.Level)
	} else {
		return fmt.Errorf("service config isn't specified")
	}

	if err := s.Cluster.Init(srvCfg.Cluster.Size, cfg.Nats.Url, nil); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {

	// start cluster
	if err := s.Cluster.Start(); err != nil {
		return err
	}

	// serve gRPC
	s.grpc.ListenAsync(ctx)

	return nil
}

func (s *serviceImpl) Close(context.Context) {
	s.Cluster.Close()
	s.grpc.Close()
}
