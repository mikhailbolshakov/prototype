package chat

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"log"
)

type Service interface {
	CreateUser(rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	CreateClientChannel(rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error)
	GetUsersStatuses(rq *pb.GetUsersStatusesRequest) (*pb.GetUserStatusesResponse, error)
}

type serviceImpl struct {
	pb.UsersClient
	pb.ChannelsClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) CreateUser(rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rs, err := u.UsersClient.CreateUser(ctx, &pb.CreateUserRequest{
		Username: rq.Username,
		Email:    rq.Email,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return rs, err
}

func (u *serviceImpl) CreateClientChannel(rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rs, err := u.ChannelsClient.CreateClientChannel(ctx, &pb.CreateClientChannelRequest{
		ClientUserId: rq.ClientUserId,
		Name:         rq.Name,
		DisplayName:  rq.DisplayName,
		Subscribers:  rq.Subscribers,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return rs, err
}

func (u *serviceImpl) GetUsersStatuses(rq *pb.GetUsersStatusesRequest) (*pb.GetUserStatusesResponse, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rs, err := u.UsersClient.GetUsersStatuses(ctx, rq)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return rs, err
}
