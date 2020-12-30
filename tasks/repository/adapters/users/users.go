package users

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type Service interface {
	GetByUserName(username string) *pb.User
	Search(request *pb.SearchRequest) (*pb.SearchResponse, error)
}

type serviceImpl struct {
	pb.UsersClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) GetByUserName(username string) *pb.User {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	user, _ := u.GetByUsername(ctx, &pb.GetByUsernameRequest{Username: username})
	return user
}

func (u *serviceImpl) Search(rq *pb.SearchRequest) (*pb.SearchResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return u.UsersClient.Search(ctx, rq)
}