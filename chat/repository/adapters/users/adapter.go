package users

import (
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/chat/domain"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
	GetService() domain.UserService
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

	usersCfg := c.Services["users"]

	cl, err := kitGrpc.NewClient(usersCfg.Grpc.Host, usersCfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.userServiceImpl.UsersClient = pb.NewUsersClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() domain.UserService {
	return a.userServiceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}