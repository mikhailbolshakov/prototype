package users

import (
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"gitlab.medzdrav.ru/prototype/users/grpc"
	"gitlab.medzdrav.ru/prototype/users/infrastructure"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
)

type serviceImpl struct {
	domain       domain.UserService
	search       domain.UserSearchService
	grpc         *grpc.Server
	mmAdapter    mattermost.Adapter
	tasksAdapter tasks.Adapter
	storage      storage.UserStorage
	infr         *infrastructure.Container
}

func New() service.Service {

	s := &serviceImpl{}

	s.infr = infrastructure.New()
	s.storage = storage.NewStorage(s.infr)
	s.mmAdapter = mattermost.NewAdapter()
	s.tasksAdapter = tasks.NewAdapter()
	s.search = domain.NewUserSearchService(s.storage)
	s.domain = domain.NewUserService(s.storage, s.mmAdapter.GetService(), s.tasksAdapter.GetService())
	s.grpc = grpc.New(s.domain, s.search)

	return s
}

func (s *serviceImpl) Init() error {

	if err := s.infr.Init(); err != nil {
		return err
	}

	if err := s.mmAdapter.Init(); err != nil {
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

	s.grpc.ListenAsync()

	if err := s.mmAdapter.ListenAsync(); err != nil {
		return err
	}

	return nil
}
