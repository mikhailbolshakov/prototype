package domain

import (
	"errors"
	"github.com/sherifabdlnaby/configuro"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	"sync"
)

type ConfigService interface {
	Load() error
	Get() (*kit.Config, error)
	GrpcSettings() *kit.Grpc
}

type serviceImpl struct {
	config *kit.Config
	sync.RWMutex
}

func New() ConfigService {
	return &serviceImpl{}
}

func (s *serviceImpl) Load() error {

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

func (s *serviceImpl) Get() (*kit.Config, error) {

	s.RLock()
	defer s.RUnlock()

	if s.config == nil {
		return nil, errors.New("config isn't loaded")
	}

	return s.config, nil

}

func (s *serviceImpl) GrpcSettings() *kit.Grpc {

	s.RLock()
	defer s.RUnlock()

	return &kit.Grpc{
		Port:  s.config.Services["config"].Grpc.Port,
		Hosts: s.config.Services["config"].Grpc.Hosts,
	}
}
