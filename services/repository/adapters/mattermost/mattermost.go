package mattermost

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/mm"
	"log"
)

type Service interface {
	CreateClientChannel(rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error)
}

type serviceImpl struct {
	pb.ChannelsClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
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

