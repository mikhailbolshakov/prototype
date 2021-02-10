package chat

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/chat/grpc"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/chat/domain/impl"
	"math/rand"
)

type serviceImpl struct {
	domainService     domain.Service
	grpc              *grpc.Server
	mattermostAdapter mattermost.Adapter
	configAdapter     config.Adapter
	configService     domain.ConfigService
	queue             queue.Queue
	queueListener     listener.QueueListener
}

func New() service.Service {

	s := &serviceImpl{}

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = &stan.Stan{}
	s.queueListener = listener.NewQueueListener(s.queue)

	s.mattermostAdapter = mattermost.NewAdapter()
	mattermostService := s.mattermostAdapter.GetService()

	s.domainService = impl.NewService(mattermostService)
	s.grpc = grpc.New(s.domainService)

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

	if err := s.mattermostAdapter.Init(c); err != nil {
		return err
	}

	if err := s.grpc.Init(c); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("mm_%d", rand.Intn(99999))); err != nil {
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
	s.configAdapter.Close()
	s.mattermostAdapter.Close()
	s.grpc.Close()
	_ = s.queue.Close()
}
