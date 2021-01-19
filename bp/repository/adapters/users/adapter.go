package users

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type Adapter interface {
	Init() error
	GetService() Service
}

type adapterImpl struct {
	userServiceImpl *serviceImpl
	initialized bool
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		userServiceImpl: newImpl(),
		initialized: false,
	}
	return a
}

func (a *adapterImpl) Init() error {
	cl, err := kitGrpc.NewClient("localhost", "50051")
	if err != nil {
		return err
	}
	a.userServiceImpl.UsersClient = pb.NewUsersClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() Service {
	return a.userServiceImpl
}
