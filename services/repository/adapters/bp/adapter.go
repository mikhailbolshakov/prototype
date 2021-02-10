package bp

import (
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/bp"
	"gitlab.medzdrav.ru/prototype/services/domain"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
	GetService() domain.BpService
	Close()
}

type adapterImpl struct {
	serviceImpl *serviceImpl
	client *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		serviceImpl: newServiceImpl(),
	}
	return a
}

func (a *adapterImpl) Init(c *kitConfig.Config) error {
	cfg := c.Services["bp"]
	cl, err := kitGrpc.NewClient(cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.serviceImpl.ProcessClient = pb.NewProcessClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() domain.BpService {
	return a.serviceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}

