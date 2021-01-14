package tasks

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"gitlab.medzdrav.ru/prototype/tasks/grpc"
	"gitlab.medzdrav.ru/prototype/tasks/infrastructure"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
	"math/rand"
)

type serviceImpl struct {
	domainTaskService       domain.TaskService
	domainConfigService     domain.ConfigService
	domainTaskSearchService domain.TaskSearchService
	assignTasksDaemon       domain.AssignmentDaemon
	grpc                    *grpc.Server
	usersAdapter            users.Adapter
	storage                 storage.TaskStorage
	infr                    *infrastructure.Container
	queue                   queue.Queue
}

func New() service.Service {

	s := &serviceImpl{}
	s.infr = infrastructure.New()
	s.storage = storage.NewStorage(s.infr)
	s.usersAdapter = users.NewAdapter()
	s.domainConfigService = domain.NewTaskConfigService()

	s.queue = &stan.Stan{}

	userService := s.usersAdapter.GetUserService()

	s.domainTaskService = domain.NewTaskService(s.storage, s.domainConfigService, userService, s.queue)
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

	if err := s.infr.Init(); err != nil {
		return err
	}

	if err := s.usersAdapter.Init(); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("client_tasks_%d", rand.Intn(99999))); err != nil {
		return err
	}

	if err := s.assignTasksDaemon.AddTaskType(&domain.Type{
		Type:    domain.TT_CLIENT,
		SubType: domain.TST_MED_REQUEST,
	}); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) Listen() error {
	return nil
}

func (s *serviceImpl) ListenAsync() error {

	s.grpc.ListenAsync()
	s.assignTasksDaemon.Run()

	return nil
}
