package config

import (
	"gitlab.medzdrav.ru/prototype/config/domain"
	"gitlab.medzdrav.ru/prototype/config/grpc"
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

func (s *serviceImpl) Init() error {
	return s.domain.Load()
}

func (s *serviceImpl) ListenAsync() error {
	s.grpc.ListenAsync()
	return nil
}

func (s *serviceImpl) Close() {
	s.grpc.Close()
}
