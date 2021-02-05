package users

import (
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
	GetService() Service
	Close()
}

type adapterImpl struct {
	userServiceImpl *serviceImpl
	client          *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		userServiceImpl: newImpl(),
	}
	return a
}

func (a *adapterImpl) Init(c *kitConfig.Config) error {
	cfg := c.Services["users"]
	cl, err := kitGrpc.NewClient(cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.userServiceImpl.UsersClient = pb.NewUsersClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() Service {
	return a.userServiceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}
