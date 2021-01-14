package mattermost

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/mm"
)

type Adapter interface {
	Init() error
	GetService() Service
}

type adapterImpl struct {
	mmServiceImpl *serviceImpl
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
	a.mmServiceImpl.ChannelsClient = pb.NewChannelsClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() Service {
	return a.mmServiceImpl
}
