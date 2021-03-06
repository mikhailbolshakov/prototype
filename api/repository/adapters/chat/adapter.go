package chat

import (
	"gitlab.medzdrav.ru/prototype/api/public"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"gitlab.medzdrav.ru/prototype/proto/config"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() public.ChatService
	Close()
}

type adapterImpl struct {
	serviceImpl *serviceImpl
	client      *kitGrpc.Client
}

func NewAdapter(userService public.UserService) Adapter {
	a := &adapterImpl{
		serviceImpl: newImpl(userService),
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
	a.serviceImpl.UsersClient = pb.NewUsersClient(cl.Conn)
	a.serviceImpl.PostsClient = pb.NewPostsClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() public.ChatService {
	return a.serviceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}