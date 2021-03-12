package services

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"gitlab.medzdrav.ru/prototype/services/domain/impl"
	"gitlab.medzdrav.ru/prototype/services/grpc"
	"gitlab.medzdrav.ru/prototype/services/logger"
	"gitlab.medzdrav.ru/prototype/services/meta"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/bp"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
)

type serviceImpl struct {
	grpc           *grpc.Server
	storageAdapter storage.Adapter
	queue          queue.Queue
	usersAdapter   users.Adapter
	configAdapter  config.Adapter
	configService  domain.ConfigService
	bpAdapter      bp.Adapter
}

func New() service.Service {

	s := &serviceImpl{}
	s.storageAdapter = storage.NewAdapter()
	strg := s.storageAdapter.GetService()

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = stan.New(logger.LF())

	s.bpAdapter = bp.NewAdapter()
	bpService := s.bpAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	userService := s.usersAdapter.GetService()

	balanceService := impl.NewBalanceService(userService, strg, s.queue)

	deliveryService := impl.NewDeliveryService(balanceService, userService, bpService, strg, s.queue)

	s.grpc = grpc.New(balanceService, deliveryService)

	return s
}

func (s *serviceImpl) GetCode() string {
	return meta.Meta.ServiceCode()
}

func (s *serviceImpl) Init(ctx context.Context) error {

	if err := s.configAdapter.Init(true); err != nil {
		return err
	}

	cfg, err := s.configService.Get()
	if err != nil {
		return err
	}

	// set logging params
	if srvCfg, ok := cfg.Services[meta.Meta.ServiceCode()]; ok {
		logger.Logger.SetLevel(srvCfg.Log.Level)
	} else {
		return fmt.Errorf("service config isn't specified")
	}

	if err := s.storageAdapter.Init(cfg); err != nil {
		return err
	}

	if err := s.grpc.Init(cfg); err != nil {
		return err
	}

	if err := s.usersAdapter.Init(cfg); err != nil {
		return err
	}

	if err := s.bpAdapter.Init(cfg); err != nil {
		return err
	}

	if err := s.queue.Open(ctx, meta.Meta.InstanceId(), &queue.Options{
		Url:       cfg.Nats.Url,
		ClusterId: cfg.Nats.ClusterId,
	}); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {
	s.grpc.ListenAsync()
	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {
	s.bpAdapter.Close()
	s.configAdapter.Close()
	s.usersAdapter.Close()
	s.grpc.Close()
	s.storageAdapter.Close()
	_ = s.queue.Close()
}
