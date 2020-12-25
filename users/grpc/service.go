package grpc

import (
	"context"
	"gitlab.medzdrav.ru/prototype/users/domain"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type Service struct {
	domain domain.UserService
	pb.UnimplementedUsersServer
}

func New() *Service {
	s := &Service{
		domain: domain.NewUserService(),
	}
	return s
}

func (s *Service) Create(ctx context.Context, rq *pb.CreateUserRequest) (*pb.User, error) {

	user, err := s.fromPb(rq)
	if err != nil {
		return nil, err
	}

	user, err = s.domain.Create(user)
	if err != nil {
		return nil, err
	}

	rs, err := s.fromDomain(user)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (s *Service) GetByUsername(ctx context.Context, rq *pb.GetByUsernameRequest) (*pb.User, error) {

	user := s.domain.GetByUsername(rq.Username)
	rs, err := s.fromDomain(user)
	if err != nil {
		return nil, err
	}

	return rs, nil
}