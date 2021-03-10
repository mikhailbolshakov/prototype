package users

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/services/domain"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() domain.UserService
	Close()
}

type adapterImpl struct {
	userServiceImpl *serviceImpl
	initialized bool
	client *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		userServiceImpl: newImpl(),
		initialized: false,
	}
	return a
}

func (a *adapterImpl) Init(c *config.Config) error {

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