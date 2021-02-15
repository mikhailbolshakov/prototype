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

func (u *serviceImpl) CreateClientChannel(ctx context.Context, rq *pb.CreateClientChannelRequest) (string, error) {
	rs, err := u.ChannelsClient.CreateClientChannel(ctx, &pb.CreateClientChannelRequest{
		ChatUserId:  rq.ChatUserId,
		Name:        rq.Name,
		DisplayName: rq.DisplayName,
		Subscribers: rq.Subscribers,
	})
	if err != nil {
		return "", err
	}

	return rs.ChannelId, err
}

func (u *serviceImpl) Subscribe(ctx context.Context, chatUserId, channelId string) error {
	_, err := u.ChannelsClient.Subscribe(ctx, &pb.SubscribeRequest{
		ChatUserId: chatUserId,
		ChannelId:  channelId,
	})
	return err
}

func (u *serviceImpl) GetChannelsForUserAndExpert(ctx context.Context, userId, expertId string) ([]string, error) {
	rs, err := u.ChannelsClient.GetChannelsForUserAndMembers(ctx, &pb.GetChannelsForUserAndMembersRequest{
		ChatUserId:        userId,
		MemberChatUserIds: []string{expertId},
	})
	if err != nil {
		return nil, err
	}
	return rs.ChannelIds, nil
}

func (u *serviceImpl) Post(ctx context.Context, message, channelId, userId string, ephemeral bool) error {
	_, err := u.PostsClient.Post(ctx, &pb.PostRequest{Posts: []*pb.Post{{
		Message:      message,
		ToChatUserId: userId,
		ChannelId:    channelId,
		Ephemeral:    ephemeral,
		From:         &pb.From{Who: pb.From_BOT},
	}}})
	return err
}

func (u *serviceImpl) PredefinedPost(ctx context.Context, channelId, userId, code string, ephemeral bool, params map[string]interface{}) error {

	var paramsB []byte
	if params != nil {
		paramsB, _ = json.Marshal(params)
	}

	_, err := u.PostsClient.Post(ctx, &pb.PostRequest{Posts: []*pb.Post{{
		ToChatUserId: userId,
		ChannelId:    channelId,
		Ephemeral:    ephemeral,
		From:         &pb.From{Who: pb.From_BOT},
		PredefinedPost: &pb.PredefinedPost{
			Code:   code,
			Params: paramsB,
		},
	}}})
	return err
}

func (u *serviceImpl) CreateUser(ctx context.Context, rq *pb.CreateUserRequest) (string, error) {
	rs, err := u.UsersClient.CreateUser(ctx, rq)
	if err != nil {
		return "", err
	}
	return rs.ChatUserId, nil
}

func (u *serviceImpl) DeleteUser(ctx context.Context, userId string) error {
	_, err := u.UsersClient.DeleteUser(ctx, &pb.DeleteUserRequest{ChatUserId: userId})
	return err
}

func (u *serviceImpl) AskBot(ctx context.Context, rq *pb.AskBotRequest) (*pb.AskBotResponse, error) {
	return u.PostsClient.AskBot(ctx, rq)
}
