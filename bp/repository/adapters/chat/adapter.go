package chat

import (
	"gitlab.medzdrav.ru/prototype/bp/domain"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"gitlab.medzdrav.ru/prototype/proto/config"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() domain.ChatService
	Close()
}

type adapterImpl struct {
	serviceImpl *serviceImpl
	client      *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		serviceImpl: newImpl(),
	}
	return a
}

func (a *adapterImpl) Init(c *config.Config) error {

	chatCfg := c.Services["chat"]
	cl, err := kitGrpc.NewClient(chatCfg.Grpc.Host, chatCfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.serviceImpl.ChannelsClient = pb.NewChannelsClient(cl.Conn)
	a.serviceImpl.PostsClient = pb.NewPostsClient(cl.Conn)
	a.serviceImpl.UsersClient = pb.NewUsersClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() domain.ChatService {
	return a.serviceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}