package sessions

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/sessions"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() domain.SessionsService
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

	cfg := c.Services["sessions"]
	cl, err := kitGrpc.NewClient(cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.serviceImpl.SessionsClient = pb.NewSessionsClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() domain.SessionsService {
	return a.serviceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}