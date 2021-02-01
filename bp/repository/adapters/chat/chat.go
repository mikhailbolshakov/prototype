package chat

import (
	"context"
	"encoding/json"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"log"
)

type Service interface {
	CreateClientChannel(rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error)
	GetChannelsForUserAndExpert(userId, expertId string) ([]string, error)
	SendTriggerPost(postCode, userId, channelId string, params map[string]interface{}) error
	Subscribe(userId, channelId string) error
	CreateUser(rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	DeleteUser(userId string) error
	AskBot(rq *pb.AskBotRequest) (*pb.AskBotResponse, error)
	SendPostFromBot(rq *pb.SendPostFromBotRequest) (*pb.SendPostFromBotResponse, error)
}

type serviceImpl struct {
	pb.ChannelsClient
	pb.PostsClient
	pb.UsersClient
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

func (u *serviceImpl) CreateUser(rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return u.UsersClient.CreateUser(context.Background(), rq)
}

func (u *serviceImpl) DeleteUser(userId string) error {
	_, err := u.UsersClient.DeleteUser(context.Background(), &pb.DeleteUserRequest{MMUserId: userId})
	return err
}

func (u *serviceImpl) AskBot(rq *pb.AskBotRequest) (*pb.AskBotResponse, error) {
	return u.PostsClient.AskBot(context.Background(), rq)
}

func (u *serviceImpl) SendPostFromBot(rq *pb.SendPostFromBotRequest) (*pb.SendPostFromBotResponse, error) {
	return u.PostsClient.SendPostFromBot(context.Background(), rq)
}