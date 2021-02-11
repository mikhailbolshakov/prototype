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
	user, _ := u.UsersClient.Get(ctx, &pb.GetByIdRequest{Id: id})
	return user
}

func (u *serviceImpl) CreateClient(ctx context.Context, request *pb.CreateClientRequest) (*pb.User, error) {
	return u.UsersClient.CreateClient(ctx, request)
}

func (u *serviceImpl) CreateConsultant(ctx context.Context, request *pb.CreateConsultantRequest) (*pb.User, error) {
	return u.UsersClient.CreateConsultant(ctx, request)
}

func (u *serviceImpl) CreateExpert(ctx context.Context, request *pb.CreateExpertRequest) (*pb.User, error) {
	return u.UsersClient.CreateExpert(ctx, request)
}

func (u *serviceImpl) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	return u.UsersClient.Search(ctx, request)
}

