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

func (u *serviceImpl) CreateClient(request *pb.CreateClientRequest) (*pb.User, error) {
	return u.UsersClient.CreateClient(context.Background(), request)
}

func (u *serviceImpl) CreateConsultant(request *pb.CreateConsultantRequest) (*pb.User, error) {
	return u.UsersClient.CreateConsultant(context.Background(), request)
}

func (u *serviceImpl) CreateExpert(request *pb.CreateExpertRequest) (*pb.User, error) {
	return u.UsersClient.CreateExpert(context.Background(), request)
}

func (u *serviceImpl) Search(request *pb.SearchRequest) (*pb.SearchResponse, error) {
	return u.UsersClient.Search(context.Background(), request)
}

