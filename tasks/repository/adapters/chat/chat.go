package chat

import (
	"context"
	"encoding/json"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type Service interface {
	SendTriggerPost(postCode, userId, channelId string, params map[string]interface{}) error
}

type serviceImpl struct {
	pb.PostsClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
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
