package users

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"gitlab.medzdrav.ru/prototype/users/grpc"
	"gitlab.medzdrav.ru/prototype/users/infrastructure"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
	"math/rand"
)

type serviceImpl struct {
	domain        domain.UserService
	search        domain.UserSearchService
	grpc          *grpc.Server
	chatAdapter   chat.Adapter
	configAdapter config.Adapter
	configService config.Service
	storage       storage.UserStorage
	infr          *infrastructure.Container
	queue         queue.Queue
}

func New() service.Service {

	s := &serviceImpl{}

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = &stan.Stan{}
	s.infr = infrastructure.New()
	s.storage = storage.NewStorage(s.infr)
	s.chatAdapter = chat.NewAdapter()

	mmService := s.chatAdapter.GetService()

	s.search = domain.NewUserSearchService(s.storage, mmService)
	s.domain = domain.NewUserService(s.storage, s.queue)
	s.grpc = grpc.New(s.domain, s.search)

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

	if err := s.queue.Open(fmt.Sprintf("users_%d", rand.Intn(99999))); err != nil {
		return err
	}

	if err := s.grpc.Init(c); err != nil {
		return err
	}

	if err := s.infr.Init(c); err != nil {
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
	s.infr.Close()
	s.grpc.Close()
}
