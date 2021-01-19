package mattermost

import (
	"context"
	"encoding/json"
	pb "gitlab.medzdrav.ru/prototype/proto/mm"
	"log"
)

type Service interface {
	CreateClientChannel(rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error)
	GetChannelsForUserAndExpert(userId, expertId string) ([]string, error)
	SendTriggerPost(postCode, userId, channelId string, params map[string]interface{}) error
	Subscribe(userId, channelId string) error
}

type serviceImpl struct {
	pb.ChannelsClient
	pb.PostsClient
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

func (u *serviceImpl) Subscribe(userId, channelId string) error {
	_, err := u.ChannelsClient.Subscribe(context.Background(), &pb.SubscribeRequest{
		UserId:    userId,
		ChannelId: channelId,
	})
	return err
}

func (u *serviceImpl) GetChannelsForUserAndExpert(userId, expertId string) ([]string, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rs, err := u.ChannelsClient.GetChannelsForUserAndMembers(ctx, &pb.GetChannelsForUserAndMembersRequest{
		UserId:        userId,
		MemberUserIds: []string{expertId},
	})
	if err != nil {
		return nil, err
	}
	return rs.ChannelIds, nil
}

func (u *serviceImpl) SendTriggerPost(postCode, userId, channelId string, params map[string]interface{}) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	paramsB, _ := json.Marshal(params)

	_, err := u.PostsClient.SendTriggerPost(ctx, &pb.SendTriggerPostRequest{
		PostCode:  postCode,
		UserId:    userId,
		ChannelId: channelId,
		Params:    paramsB,
	})

	return err

}
