package users

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type Adapter interface {
	Init() error
	GetService() Service
	Close()
}

type adapterImpl struct {
	userServiceImpl *serviceImpl
	client *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		userServiceImpl: newImpl(),
	}
	return a
}

func (s *adapterImpl) Init() error {
	cl, err := kitGrpc.NewClient("localhost", "50051")
	if err != nil {
		return err
	}
	s.client = cl
	s.userServiceImpl.UsersClient = pb.NewUsersClient(cl.Conn)
	return nil
}

func (s *adapterImpl) GetService() Service {
	return s.userServiceImpl
}

func (s *adapterImpl) Close() {
	_ = s.client.Conn.Close()
}
