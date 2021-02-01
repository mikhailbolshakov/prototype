package bp

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/bp"
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
	s := &adapterImpl{
		serviceImpl: newImpl(),
	}
	return s
}

func (s *adapterImpl) Init() error {
	cl, err := kitGrpc.NewClient("localhost", "50055")
	if err != nil {
		return err
	}
	s.client = cl
	s.serviceImpl.ProcessClient = pb.NewProcessClient(cl.Conn)
	return nil
}

func (s *adapterImpl) GetService() Service {
	return s.serviceImpl
}

func (s *adapterImpl) Close() {
	_ = s.client.Conn.Close()
}