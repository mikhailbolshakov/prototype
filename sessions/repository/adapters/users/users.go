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

func (u *serviceImpl) Get(ctx context.Context, id string) *pb.User {
	usr, _ := u.UsersClient.Get(ctx, &pb.GetByIdRequest{Id: id})
	return usr
}
