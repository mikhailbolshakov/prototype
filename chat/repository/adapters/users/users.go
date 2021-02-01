package users

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type Service interface {
	GetByMMId(mmId string) *pb.User
	Get(id string) *pb.User
}

type serviceImpl struct {
	pb.UsersClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) GetByMMId(mmId string) *pb.User {
	user, _ := u.UsersClient.GetByMMId(context.Background(), &pb.GetByMMIdRequest{MMId: mmId})
	return user
}

func (u *serviceImpl) Get(id string) *pb.User {
	user, _ := u.UsersClient.Get(context.Background(), &pb.GetByIdRequest{Id: id})
	return user
}
