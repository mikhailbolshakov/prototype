package chat

import (
	"context"
	"encoding/json"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type serviceImpl struct {
	pb.PostsClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
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

	_, err := u.PostsClient.Post(ctx, &pb.PostRequest{Posts: []*pb.Post{&pb.Post{
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
