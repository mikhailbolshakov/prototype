package tasks

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"gitlab.medzdrav.ru/prototype/tasks/grpc"
	"gitlab.medzdrav.ru/prototype/tasks/infrastructure"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
	"math/rand"
)

type serviceImpl struct {
	domainTaskService       domain.TaskService
	domainConfigService     domain.ConfigService
	domainTaskSearchService domain.TaskSearchService
	assignTasksDaemon       domain.AssignmentDaemon
	configAdapter           config.Adapter
	configService           config.Service
	scheduler               domain.TaskScheduler
	grpc                    *grpc.Server
	usersAdapter            users.Adapter
	chatAdapter             chat.Adapter
	storage                 storage.TaskStorage
	infr                    *infrastructure.Container
	queue                   queue.Queue
}

func New() service.Service {

	s := &serviceImpl{}


	s.queue = &stan.Stan{}
	s.infr = infrastructure.New()
	s.storage = storage.NewStorage(s.infr)
	s.domainConfigService = domain.NewTaskConfigService()
	s.scheduler = domain.NewScheduler(s.domainConfigService, s.storage)

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	userService := s.usersAdapter.GetUserService()

	s.chatAdapter = chat.NewAdapter()
	chatService := s.chatAdapter.GetService()

	s.domainTaskService = domain.NewTaskService(s.scheduler, s.storage, s.domainConfigService, userService, s.queue, chatService)
	s.domainTaskSearchService = domain.NewTaskSearchService(s.storage)

	s.assignTasksDaemon = domain.NewAssignmentDaemon(s.domainConfigService,
		s.domainTaskService,
		s.domainTaskSearchService,
		userService,
		s.storage)

	s.grpc = grpc.New(s.domainTaskService, s.domainTaskSearchService)

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

	if err := s.chatAdapter.Init(c); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("client_tasks_%d", rand.Intn(99999))); err != nil {
		return err
	}

	if err := s.assignTasksDaemon.Init(); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync() error {

	s.grpc.ListenAsync()
	s.assignTasksDaemon.Run()
	s.scheduler.StartAsync()

	return nil
}

func (s *serviceImpl) Close() {

	s.configAdapter.Close()

	// TODO: if uncomment hangs on, has to be investigated
	_ = s.assignTasksDaemon.Stop()
	s.chatAdapter.Close()
	s.usersAdapter.Close()

	_ = s.queue.Close()
	s.infr.Close()
	s.grpc.Close()
}
