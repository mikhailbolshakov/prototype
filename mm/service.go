package mm

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/mm/domain"
	"gitlab.medzdrav.ru/prototype/mm/grpc"
	"gitlab.medzdrav.ru/prototype/mm/infrastructure"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/users"
	"math/rand"
)

type serviceImpl struct {
	domainService     domain.Service
	grpc              *grpc.Server
	mattermostAdapter mattermost.Adapter
	usersAdapter      users.Adapter
	tasksAdapter      tasks.Adapter
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

	s.domainService = domain.NewService(mattermostService, usersService, tasksService, s.infr.Bpm)
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

	s.queueListener.Add("tasks.remind", s.domainService.TaskRemindMessageHandler)
	s.queueListener.Add("tasks.duedate", s.domainService.TaskDueDateMessageHandler)
	s.queueListener.Add("mm.posts", s.domainService.MattermostPostMessageHandler)

	return nil

}

func (s *serviceImpl) Listen() error {
	return nil
}

func (s *serviceImpl) ListenAsync() error {

	s.grpc.ListenAsync()
	s.queueListener.ListenAsync()

	return nil
}
