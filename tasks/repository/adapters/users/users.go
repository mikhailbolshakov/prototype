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

func (u *serviceImpl) Get(ctx context.Context, id, username string) *pb.User {

	if id != "" {
		user, _ := u.UsersClient.Get(ctx, &pb.GetByIdRequest{Id: id})
		return user
	} else if username != "" {
		user, _ := u.UsersClient.Get(ctx, &pb.GetByIdRequest{Id: username})
		return user
	}
	return nil
}

func (u *serviceImpl) Search(ctx context.Context, rq *pb.SearchRequest) (*pb.SearchResponse, error) {
	return u.UsersClient.Search(ctx, rq)
}