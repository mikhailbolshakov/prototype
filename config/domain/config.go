package domain

import (
	"context"
	"errors"
	"github.com/sherifabdlnaby/configuro"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	"sync"
)

type ConfigService interface {
	Load(ctx context.Context) error
	Get(ctx context.Context) (*kit.Config, error)
	GrpcSettings(ctx context.Context) *kit.Grpc
}

type serviceImpl struct {
	config *kit.Config
	sync.RWMutex
}

func New() ConfigService {
	return &serviceImpl{}
}

func (s *serviceImpl) Load(context.Context) error {

	s.Lock()
	defer s.Unlock()

	Loader, err := configuro.NewConfig(
		configuro.WithLoadFromConfigFile("../config.yml", true),
		configuro.WithLoadDotEnv("../.env"))
	if err != nil {
		return err
	}

	s.config = &kit.Config{}

	err = Loader.Load(s.config)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) Get(context.Context) (*kit.Config, error) {

	s.RLock()
	defer s.RUnlock()

	if s.config == nil {
		return nil, errors.New("config isn't loaded")
	}

	return s.config, nil

}

func (s *serviceImpl) GrpcSettings(context.Context) *kit.Grpc {

	s.RLock()
	defer s.RUnlock()

	return &kit.Grpc{
		Port: s.config.Services["cfg"].Grpc.Port,
		Host: s.config.Services["cfg"].Grpc.Host,
	}
}
