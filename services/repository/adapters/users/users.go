package users

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type serviceImpl struct {
	pb.UsersClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) Get(id string) *pb.User {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	user, _ := u.UsersClient.Get(ctx, &pb.GetByIdRequest{Id: id})
	return user
}

