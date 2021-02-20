package chat

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type serviceImpl struct {
	pb.UsersClient
	pb.ChannelsClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) CreateUser(ctx context.Context, rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	rs, err := u.UsersClient.CreateUser(ctx, &pb.CreateUserRequest{
		Username: rq.Username,
		Email:    rq.Email,
	})
	if err != nil {
		return nil, err
	}

	return rs, err
}

func (u *serviceImpl) CreateClientChannel(ctx context.Context, rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error) {
	rs, err := u.ChannelsClient.CreateClientChannel(ctx, &pb.CreateClientChannelRequest{
		ChatUserId:  rq.ChatUserId,
		Name:        rq.Name,
		DisplayName: rq.DisplayName,
		Subscribers: rq.Subscribers,
	})
	if err != nil {
		return nil, err
	}

	return rs, err
}

func (u *serviceImpl) GetUsersStatuses(ctx context.Context, rq *pb.GetUsersStatusesRequest) (*pb.GetUserStatusesResponse, error) {

	rs, err := u.UsersClient.GetUsersStatuses(ctx, rq)
	if err != nil {
		return nil, err
	}

	return rs, err
}
