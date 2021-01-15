package services

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"gitlab.medzdrav.ru/prototype/services/grpc"
	"gitlab.medzdrav.ru/prototype/services/infrastructure"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
	"math/rand"
)

type serviceImpl struct {
	balanceService  domain.UserBalanceService
	deliveryService domain.DeliveryService
	tasksAdapter    tasks.Adapter
	usersAdapter    users.Adapter
	mmAdapter       mattermost.Adapter
	grpc            *grpc.Server
	storage         storage.Storage
	infr            *infrastructure.Container
	queue           queue.Queue
	bpm  			bpm.Engine
}

func New() service.Service {

	s := &serviceImpl{}
	s.infr = infrastructure.New()
	s.storage = storage.NewStorage(s.infr)

	s.queue = &stan.Stan{}
	s.bpm = s.infr.Bpm

	s.tasksAdapter = tasks.NewAdapter(s.queue)
	taskService := s.tasksAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	userService := s.usersAdapter.GetService()

	s.mmAdapter = mattermost.NewAdapter()
	mmService := s.mmAdapter.GetService()

	s.balanceService = domain.NewBalanceService(s.storage, s.queue)
	s.deliveryService = domain.NewDeliveryService(s.balanceService, taskService, userService, mmService, s.storage, s.queue, s.bpm)

	s.grpc = grpc.New(s.balanceService, s.deliveryService)

	return s
}

func (s *serviceImpl) Init() error {

	if err := s.infr.Init(); err != nil {
		return err
	}

	if err := s.tasksAdapter.Init(); err != nil {
		return err
	}

	if err := s.usersAdapter.Init(); err != nil {
		return err
	}

	if err := s.mmAdapter.Init(); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("client_tasks_%d", rand.Intn(99999))); err != nil {
		return err
	}

	if err := s.bpm.DeployBPMNs([]string{"../services/bpmn/expert_online_consultation.bpmn"}); err != nil {
		return err
	}

	if err := s.deliveryService.RegisterBpmHandlers(); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) Listen() error {
	return nil
}

func (s *serviceImpl) ListenAsync() error {

	s.grpc.ListenAsync()
	s.tasksAdapter.ListenAsync()

	return nil
}
