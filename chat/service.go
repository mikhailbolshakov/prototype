package chat

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/chat/grpc"
	"gitlab.medzdrav.ru/prototype/chat/infrastructure"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/bp"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/users"
	"math/rand"
)

type serviceImpl struct {
	domainService     domain.Service
	grpc              *grpc.Server
	mattermostAdapter mattermost.Adapter
	usersAdapter      users.Adapter
	tasksAdapter      tasks.Adapter
	bpAdapter         bp.Adapter
	queue             queue.Queue
	queueListener     listener.QueueListener
	infr              *infrastructure.Container
}

func New() service.Service {

	s := &serviceImpl{}

	s.infr = infrastructure.New()

	s.queue = &stan.Stan{}
	s.queueListener = listener.NewQueueListener(s.queue)

	s.mattermostAdapter = mattermost.NewAdapter(s.queue)
	mattermostService := s.mattermostAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	usersService := s.usersAdapter.GetService()

	s.tasksAdapter = tasks.NewAdapter(s.queue)
	tasksService := s.tasksAdapter.GetService()

	s.bpAdapter = bp.NewAdapter()
	bpService := s.bpAdapter.GetService()

	s.domainService = domain.NewService(mattermostService, usersService, tasksService, bpService)
	s.grpc = grpc.New(s.domainService)

	return s
}

func (s *serviceImpl) Init() error {

	if err := s.infr.Init(); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("mm_%d", rand.Intn(99999))); err != nil {
		return err
	}

	if err := s.mattermostAdapter.Init(); err != nil {
		return err
	}

	if err := s.usersAdapter.Init(); err != nil {
		return err
	}

	if err := s.tasksAdapter.Init(); err != nil {
		return err
	}

	if err := s.bpAdapter.Init(); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync() error {

	s.grpc.ListenAsync()
	s.queueListener.ListenAsync()

	return nil
}

func (s *serviceImpl) Close() {
	s.mattermostAdapter.Close()
	s.tasksAdapter.Close()
	s.usersAdapter.Close()
	s.bpAdapter.Close()
	s.grpc.Close()
	s.infr.Close()
	_ = s.queue.Close()
}
