package tasks

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"gitlab.medzdrav.ru/prototype/tasks/domain/impl"
	"gitlab.medzdrav.ru/prototype/tasks/grpc"
	"gitlab.medzdrav.ru/prototype/tasks/logger"
	"gitlab.medzdrav.ru/prototype/tasks/meta"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
)

type serviceImpl struct {
	taskService       domain.TaskService
	taskConfigService domain.ConfigService
	assignTasksDaemon domain.AssignmentDaemon
	configAdapter     config.Adapter
	cfgService        domain.CfgService
	scheduler         domain.TaskScheduler
	grpc              *grpc.Server
	usersAdapter      users.Adapter
	chatAdapter       chat.Adapter
	storageAdapter    storage.Adapter
	queue             queue.Queue
}

func New() service.Service {

	s := &serviceImpl{}

	s.queue = stan.New(logger.LF())
	s.storageAdapter = storage.NewAdapter()
	strg := s.storageAdapter.GetService()
	s.taskConfigService = impl.NewTaskConfigService()
	s.scheduler = impl.NewScheduler(s.taskConfigService, strg)

	s.configAdapter = config.NewAdapter()
	s.cfgService = s.configAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	userService := s.usersAdapter.GetUserService()

	s.chatAdapter = chat.NewAdapter()
	chatService := s.chatAdapter.GetService()

	s.taskService = impl.NewTaskService(s.scheduler, strg, s.taskConfigService, userService, s.queue, chatService)

	s.assignTasksDaemon = impl.NewAssignmentDaemon(s.taskConfigService,
		s.taskService,
		userService,
		strg)

	s.grpc = grpc.New(s.taskService)

	return s
}

func (s *serviceImpl) GetCode() string {
	return meta.ServiceCode
}

func (s *serviceImpl) Init(ctx context.Context) error {

	if err := s.configAdapter.Init(true); err != nil {
		return err
	}

	cfg, err := s.cfgService.Get(ctx)
	if err != nil {
		return err
	}

	// set logging params
	if srvCfg, ok := cfg.Services[meta.ServiceCode]; ok {
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

	if err := s.chatAdapter.Init(cfg); err != nil {
		return err
	}

	if err := s.queue.Open(ctx, meta.NodeId, &queue.Options{
		Url:       cfg.Nats.Url,
		ClusterId: cfg.Nats.ClusterId,
	}); err != nil {
		return err
	}

	if err := s.assignTasksDaemon.Init(ctx); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {

	s.grpc.ListenAsync()
	//s.assignTasksDaemon.Run(ctx)
	//s.scheduler.StartAsync(ctx)

	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {

	s.configAdapter.Close()

	_ = s.assignTasksDaemon.Stop(ctx)
	s.chatAdapter.Close()
	s.usersAdapter.Close()

	_ = s.queue.Close()
	s.storageAdapter.Close()
	s.grpc.Close()
}
