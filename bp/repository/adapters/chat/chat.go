package chat

import (
	"context"
	"encoding/json"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type serviceImpl struct {
	pb.ChannelsClient
	pb.PostsClient
	pb.UsersClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) CreateClientChannel(rq *pb.CreateClientChannelRequest) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rs, err := u.ChannelsClient.CreateClientChannel(ctx, &pb.CreateClientChannelRequest{
		ClientUserId: rq.ClientUserId,
		Name:         rq.Name,
		DisplayName:  rq.DisplayName,
		Subscribers:  rq.Subscribers,
	})
	if err != nil {
		return "", err
	}

	return rs.ChannelId, err
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

func (u *serviceImpl) Post(message, channelId, userId string, ephemeral, fromBot bool) error {
	_, err := u.PostsClient.Post(context.Background(), &pb.PostRequest{Posts: []*pb.Post{&pb.Post{
		Message:        message,
		ToUserId:       userId,
		ChannelId:      channelId,
		Ephemeral:      ephemeral,
		FromBot:        fromBot,
	}}})
	return err
}

func (u *serviceImpl) PredefinedPost(channelId, userId, code string, ephemeral, fromBot bool, params map[string]interface{}) error {

	var paramsB []byte
	if params != nil {
		paramsB, _ = json.Marshal(params)
	}

	_, err := u.PostsClient.Post(context.Background(), &pb.PostRequest{Posts: []*pb.Post{&pb.Post{
		ToUserId:       userId,
		ChannelId:      channelId,
		Ephemeral:      ephemeral,
		FromBot:        fromBot,
		PredefinedPost: &pb.PredefinedPost{
			Code:   code,
			Params: paramsB,
		},
	}}})
	return err
}

func (u *serviceImpl) CreateUser(rq *pb.CreateUserRequest) (string, error) {
	rs, err := u.UsersClient.CreateUser(context.Background(), rq)
	if err != nil {
		return "", err
	}
	return rs.Id, nil
}

func (u *serviceImpl) DeleteUser(userId string) error {
	_, err := u.UsersClient.DeleteUser(context.Background(), &pb.DeleteUserRequest{MMUserId: userId})
	return err
}

func (u *serviceImpl) AskBot(rq *pb.AskBotRequest) (*pb.AskBotResponse, error) {
	return u.PostsClient.AskBot(context.Background(), rq)
}
