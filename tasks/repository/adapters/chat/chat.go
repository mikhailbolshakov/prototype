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

func (u *serviceImpl) Post(ctx context.Context, message, channelId, userId string, ephemeral, fromBot bool) error {
	_, err := u.PostsClient.Post(ctx, &pb.PostRequest{Posts: []*pb.Post{&pb.Post{
		Message:        message,
		ToUserId:       userId,
		ChannelId:      channelId,
		Ephemeral:      ephemeral,
		FromBot:        fromBot,
	}}})
	return err
}

func (u *serviceImpl) PredefinedPost(ctx context.Context, channelId, userId, code string, ephemeral, fromBot bool, params map[string]interface{}) error {

	var paramsB []byte
	if params != nil {
		paramsB, _ = json.Marshal(params)
	}

	_, err := u.PostsClient.Post(ctx, &pb.PostRequest{Posts: []*pb.Post{&pb.Post{
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
