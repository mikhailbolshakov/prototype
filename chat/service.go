package chat

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/chat/domain/impl"
	"gitlab.medzdrav.ru/prototype/chat/grpc"
	"gitlab.medzdrav.ru/prototype/chat/logger"
	"gitlab.medzdrav.ru/prototype/chat/meta"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
)

type serviceImpl struct {
	domainService     domain.Service
	grpc              *grpc.Server
	mattermostAdapter mattermost.Adapter
	configAdapter     config.Adapter
	configService     domain.ConfigService
	queue             queue.Queue
	queueListener     listener.QueueListener
	chatSessionHub    mattermost.ChatSessionHub
	userService       domain.UserService
	userAdapter       users.Adapter
}

func New() service.Service {

	s := &serviceImpl{}

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.userAdapter = users.NewAdapter()
	s.userService = s.userAdapter.GetService()

	s.queue = stan.New(logger.LF())
	s.queueListener = listener.NewQueueListener(s.queue, logger.LF())

	return s
}

func (s *serviceImpl) GetCode() string {
	return meta.ServiceCode
}

func (s *serviceImpl) Init(ctx context.Context) error {

	if err := s.configAdapter.Init(true); err != nil {
		return err
	}

	cfg, err := s.configService.Get()
	if err != nil {
		return err
	}

	// set logging params
	if srvCfg, ok := cfg.Services[meta.ServiceCode]; ok {
		logger.Logger.SetLevel(srvCfg.Log.Level)
	} else {
		return fmt.Errorf("service config isn't specified")
	}

	s.chatSessionHub = mattermost.NewHub(cfg, s.userService, s.queue)

	s.mattermostAdapter = mattermost.NewAdapter(s.chatSessionHub)
	mattermostService := s.mattermostAdapter.GetService()

	s.domainService = impl.NewService(mattermostService)
	s.grpc = grpc.New(s.domainService)

	if err := s.mattermostAdapter.Init(ctx, cfg); err != nil {
		return err
	}

	if err := s.userAdapter.Init(cfg); err != nil {
		return err
	}

	if err := s.grpc.Init(cfg); err != nil {
		return err
	}

	if err := s.queue.Open(ctx, meta.NodeId, &queue.Options{
		Url:       cfg.Nats.Url,
		ClusterId: cfg.Nats.ClusterId,
	}); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {

	s.grpc.ListenAsync()
	s.queueListener.ListenAsync()

	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {
	s.configAdapter.Close()
	s.userAdapter.Close()
	s.mattermostAdapter.Close(ctx)
	s.grpc.Close()
	_ = s.queue.Close()
}
