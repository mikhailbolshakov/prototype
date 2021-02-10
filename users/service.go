package users

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"gitlab.medzdrav.ru/prototype/users/domain/impl"
	"gitlab.medzdrav.ru/prototype/users/grpc"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
	"math/rand"
)

type serviceImpl struct {
	domainService  domain.UserService
	grpc           *grpc.Server
	chatAdapter    chat.Adapter
	configAdapter  config.Adapter
	configService  domain.ConfigService
	storageAdapter storage.Adapter
	queue          queue.Queue
}

func New() service.Service {

	s := &serviceImpl{}

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = &stan.Stan{}
	s.storageAdapter = storage.NewAdapter()
	strg := s.storageAdapter.GetService()
	s.chatAdapter = chat.NewAdapter()

	chatService := s.chatAdapter.GetService()

	s.domainService = impl.NewUserService(strg, chatService, s.queue)
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

	if err := s.storageAdapter.Init(c); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("users_%d", rand.Intn(99999))); err != nil {
		return err
	}

	if err := s.grpc.Init(c); err != nil {
		return err
	}

	if err := s.chatAdapter.Init(c); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync() error {
	s.grpc.ListenAsync()
	return nil
}

func (s *serviceImpl) Close() {
	s.configAdapter.Close()
	s.chatAdapter.Close()
	_ = s.queue.Close()
	s.storageAdapter.Close()
	s.grpc.Close()
}
