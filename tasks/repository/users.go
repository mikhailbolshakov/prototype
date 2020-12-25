package repository

import (
	"context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type UsersServiceAdapter interface {
	GetByUserName(username string) *pb.User
}

type UsersServiceAdapterImpl struct {
	pb.UsersClient
}

func NewUsersServiceAdapter() UsersServiceAdapter {
	a := &UsersServiceAdapterImpl{}
	cl, _ := kitGrpc.NewClient("localhost", "50051")
	a.UsersClient = pb.NewUsersClient(cl.Conn)
	return a
}

func (u *UsersServiceAdapterImpl) GetByUserName(username string) *pb.User {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	user, _ := u.GetByUsername(ctx, &pb.GetByUsernameRequest{Username: username})
	return user
}