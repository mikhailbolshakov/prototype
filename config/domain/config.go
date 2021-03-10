package domain

import (
	"context"
	"fmt"
	"github.com/sherifabdlnaby/configuro"
	"gitlab.medzdrav.ru/prototype/config/logger"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"os"
	"path/filepath"
	"sync"
)

type ConfigService interface {
	Load(ctx context.Context) error
	Get(ctx context.Context) (*config.Config, error)
	GrpcSettings(ctx context.Context) *config.Grpc
}

type serviceImpl struct {
	config *config.Config
	sync.RWMutex
}

func New() ConfigService {
	return &serviceImpl{}
}

func (s *serviceImpl) Load(context.Context) error {

	l := logger.L().Cmp("config").Mth("load")

	configPath := os.Getenv("CONFIG_SOURCE_PATH")
	if configPath == "" {
		err := fmt.Errorf("env var CONFIG_SOURCE_PATH is empty")
		l.E(err).St().Err()
		return err
	}

	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			err := fmt.Errorf("config file %s not found", configPath)
			l.E(err).St().Err()
			return err
		}
	}

	absPath, _ := filepath.Abs(configPath)
	l.Dbg("loading config from file", absPath)

	s.Lock()
	defer s.Unlock()

	Loader, err := configuro.NewConfig(
		configuro.WithLoadFromConfigFile(configPath, true),
		configuro.WithLoadDotEnv(".env"))
	if err != nil {
		l.E(err).St().Err()
		return err
	}

	s.config = &config.Config{}

	err = Loader.Load(s.config)
	if err != nil {
		l.E(err).St().Err()
		return err
	}

	l.TrcF(kit.ToJson(s.config))

	return nil
}

func (s *serviceImpl) Get(context.Context) (*config.Config, error) {

	l := logger.L().Cmp("config").Mth("get").Dbg()

	s.RLock()
	defer s.RUnlock()

	if s.config == nil {
		err := fmt.Errorf("config isn't loaded")
		l.E(err).St().Err()
		return nil, err
	}

	return s.config, nil

}

func (s *serviceImpl) GrpcSettings(context.Context) *config.Grpc {

	s.RLock()
	defer s.RUnlock()

	return &config.Grpc{
		Port: s.config.Services["cfg"].Grpc.Port,
		Host: s.config.Services["cfg"].Grpc.Host,
	}
}
