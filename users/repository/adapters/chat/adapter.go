package chat

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type Adapter interface {
	Init() error
	GetService() Service
	Close()
}

type adapterImpl struct {
	mmServiceImpl *serviceImpl
	client *kitGrpc.Client
	initialized bool
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		mmServiceImpl: newImpl(),
		initialized:   false,
	}
	return a
}

func (a *adapterImpl) Init() error {

	cl, err := kitGrpc.NewClient("localhost", "50053")
	if err != nil {
		return err
	}
	a.client = cl
	a.mmServiceImpl.UsersClient = pb.NewUsersClient(cl.Conn)
	a.mmServiceImpl.ChannelsClient = pb.NewChannelsClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() Service {
	return a.mmServiceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}