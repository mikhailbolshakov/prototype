package webrtc

import (
	"gitlab.medzdrav.ru/prototype/api/public"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/webrtc"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() public.RoomService
	Close()
}

type adapterImpl struct {
	roomServiceImpl *serviceImpl
	client *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		roomServiceImpl: newImpl(),
	}
	return a
}

func (a *adapterImpl) Init(c *config.Config) error {
	cfg := c.Services["webrtc"]
	cl, err := kitGrpc.NewClient(cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.roomServiceImpl.RoomsClient = pb.NewRoomsClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() public.RoomService {
	return a.roomServiceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}

