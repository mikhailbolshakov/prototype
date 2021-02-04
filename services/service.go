package services

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"gitlab.medzdrav.ru/prototype/services/grpc"
	"gitlab.medzdrav.ru/prototype/services/infrastructure"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/bp"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
	"math/rand"
)

type serviceImpl struct {
	grpc          *grpc.Server
	storage       storage.Storage
	infr          *infrastructure.Container
	queue         queue.Queue
	usersAdapter  users.Adapter
	configAdapter config.Adapter
	configService config.Service
	bpAdapter     bp.Adapter
}

func New() service.Service {

	s := &serviceImpl{}
	s.infr = infrastructure.New()
	s.storage = storage.NewStorage(s.infr)

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = &stan.Stan{}

	s.bpAdapter = bp.NewAdapter()
	bpService := s.bpAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	userService := s.usersAdapter.GetService()

	balanceService := domain.NewBalanceService(userService, s.storage, s.queue)

	deliveryService := domain.NewDeliveryService(balanceService, userService, bpService, s.storage, s.queue)

	s.grpc = grpc.New(balanceService, deliveryService)

	return s
}

func (s *serviceImpl) Init() error {

	if err := s.configAdapter.Init(); err != nil {
		return err
	}

	c, err := s.configService.Get()
	if err != nil {
		return err
	}

	if err := s.infr.Init(c); err != nil {
		return err
	}

	if err := s.grpc.Init(c); err != nil {
		return err
	}

	if err := s.usersAdapter.Init(c); err != nil {
		return err
	}

	if err := s.bpAdapter.Init(c); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("client_tasks_%d", rand.Intn(99999))); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync() error {
	s.grpc.ListenAsync()
	return nil
}

func (s *serviceImpl) Close() {
	s.bpAdapter.Close()
	s.configAdapter.Close()
	s.usersAdapter.Close()
	s.grpc.Close()
	s.infr.Close()
	_ = s.queue.Close()
}
