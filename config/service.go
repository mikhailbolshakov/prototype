package config

import (
	"context"
	"gitlab.medzdrav.ru/prototype/config/domain"
	"gitlab.medzdrav.ru/prototype/config/grpc"
	"gitlab.medzdrav.ru/prototype/config/meta"
	"gitlab.medzdrav.ru/prototype/kit/service"
)

type serviceImpl struct {
	domain domain.ConfigService
	grpc *grpc.Server
}

func New() service.Service {
	s := &serviceImpl{}
	s.domain = domain.New()
	s.grpc = grpc.New(s.domain)
	return s
}

func (s *serviceImpl) GetCode() string {
	return meta.ServiceCode
}

func (s *serviceImpl) Init(ctx context.Context) error {
	return s.domain.Load(ctx)
}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {
	s.grpc.ListenAsync(ctx)
	return nil
}

func (s *serviceImpl) Close(context.Context) {
	s.grpc.Close()
}
