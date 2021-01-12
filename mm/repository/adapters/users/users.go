package users

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type Service interface {
	GetByMMId(mmId string) *pb.User
	GetByUsername(username string) *pb.User
}

type serviceImpl struct {
	pb.UsersClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) GetByMMId(mmId string) *pb.User {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	user, _ := u.UsersClient.GetByMMId(ctx, &pb.GetByMMIdRequest{MMId: mmId})
	return user
}

func (u *serviceImpl) GetByUsername(username string) *pb.User {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	user, _ := u.UsersClient.GetByUsername(ctx, &pb.GetByUsernameRequest{Username: username})
	return user
}
