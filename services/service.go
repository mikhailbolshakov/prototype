package services

import (
	"fmt"
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
	domainService domain.UserBalanceService
	grpc          *grpc.Server
	storage       storage.BalanceStorage
	infr          *infrastructure.Container
	queue         queue.Queue
}

func New() service.Service {

	s := &serviceImpl{}
	s.infr = infrastructure.New()
	s.storage = storage.NewStorage(s.infr)

	s.queue = &stan.Stan{}

	s.domainService = domain.NewService(s.storage, s.queue)

	s.grpc = grpc.New(s.domainService)

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
