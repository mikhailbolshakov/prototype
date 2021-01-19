package services

import (
	"fmt"
	bpmKit "gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"gitlab.medzdrav.ru/prototype/services/grpc"
	"gitlab.medzdrav.ru/prototype/services/infrastructure"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
	"math/rand"
)

type serviceImpl struct {
	grpc         *grpc.Server
	storage      storage.Storage
	infr         *infrastructure.Container
	queue        queue.Queue
	bpm          bpmKit.Engine
}

func New() service.Service {

	s := &serviceImpl{}
	s.infr = infrastructure.New()
	s.storage = storage.NewStorage(s.infr)

	s.queue = &stan.Stan{}
	s.bpm = s.infr.Bpm

	balanceService := domain.NewBalanceService(s.storage, s.queue)
	deliveryService := domain.NewDeliveryService(balanceService, s.storage, s.queue, s.bpm)

	s.grpc = grpc.New(balanceService, deliveryService)

	return s
}

func (s *serviceImpl) Init() error {

	if err := s.infr.Init(); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("client_tasks_%d", rand.Intn(99999))); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) Listen() error {
	return nil
}

func (s *serviceImpl) ListenAsync() error {
	s.grpc.ListenAsync()
	return nil
}
