package webrtc

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/webrtc"
)

type serviceImpl struct {
	pb.RoomsClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) Create(ctx context.Context, channelId string) (*pb.Room, error) {
	return u.RoomsClient.Create(ctx, &pb.CreateRoomRequest{ChannelId: channelId})
}

func (u *serviceImpl) Get(ctx context.Context, roomId string) (*pb.Room, error) {
	return u.RoomsClient.Get(ctx, &pb.GetRoomRequest{Id: roomId})
}
