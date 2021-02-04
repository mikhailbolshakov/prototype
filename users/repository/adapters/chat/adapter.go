package chat

import (
	"gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type Adapter interface {
	Init(c *config.Config) error
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

func (a *adapterImpl) Init(c *config.Config) error {

	chatCfg := c.Services["chat"]

	cl, err := kitGrpc.NewClient(chatCfg.Grpc.Hosts[0], chatCfg.Grpc.Port)
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