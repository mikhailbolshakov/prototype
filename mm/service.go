package mm

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/mm/domain"
	"gitlab.medzdrav.ru/prototype/mm/grpc"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/users"
	"math/rand"
)

type serviceImpl struct {
	domainMMService   domain.Service
	grpc              *grpc.Server
	mattermostAdapter mattermost.Adapter
	usersAdapter 	  users.Adapter
	tasksAdapter 	  tasks.Adapter
	queue             queue.Queue
}

func New() service.Service {

	s := &serviceImpl{}
	s.queue = &stan.Stan{}
	s.mattermostAdapter = mattermost.NewAdapter(s.queue)
	mattermostService := s.mattermostAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	usersService := s.usersAdapter.GetService()

	s.tasksAdapter = tasks.NewAdapter(s.queue)
	tasksService := s.tasksAdapter.GetService()

	s.domainMMService = domain.NewService(mattermostService, usersService, tasksService)
	s.grpc = grpc.New(s.domainMMService)

	return s
}

func (s *serviceImpl) Init() error {

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

	return nil

}

func (s *serviceImpl) Listen() error {
	return nil
}

func (s *serviceImpl) ListenAsync() error {

	if err := s.mattermostAdapter.ListenAsync(); err != nil {
		return err
	}

	if err := s.tasksAdapter.ListenAsync(); err != nil {
		return err
	}

	s.grpc.ListenAsync()

	return nil
}
