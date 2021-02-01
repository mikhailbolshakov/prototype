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
	serviceImpl *serviceImpl
	client *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		serviceImpl: newImpl(),
	}
	return a
}

func (a *adapterImpl) Init() error {

	cl, err := kitGrpc.NewClient("localhost", "50053")
	if err != nil {
		return err
	}
	a.client = cl
	a.serviceImpl.PostsClient = pb.NewPostsClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() Service {
	return a.serviceImpl
}

func (a *adapterImpl) Close()  {
	_ = a.client.Conn.Close()
}
