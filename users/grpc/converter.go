package grpc

import (
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/users/domain"
)

func (s *Service) fromPb(request *pb.CreateUserRequest) (*domain.User, error) {
	return &domain.User{
		Type:      request.Type,
		Username:  request.Username,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		Email:     request.Email,
	}, nil
}

func (s *Service) fromDomain(user *domain.User) (*pb.User, error) {

	if user == nil {
		return nil, nil
	}

	return &pb.User{
			Id:        user.Id,
			Type:      user.Type,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
			MMId:      user.MMUserId,
		},
		nil
}
