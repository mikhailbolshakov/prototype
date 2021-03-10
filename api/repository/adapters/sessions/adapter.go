package sessions

import (
	"gitlab.medzdrav.ru/prototype/api/public"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/sessions"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() public.SessionsService
	GetMonitor() public.SessionMonitor
	Close()
}

type adapterImpl struct {
	serviceImpl *serviceImpl
	monitorImpl *monitorImpl
	client      *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		serviceImpl: newImpl(),
		monitorImpl: newMonitorImpl(),
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
	a.monitorImpl.MonitorClient = pb.NewMonitorClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() public.SessionsService {
	return a.serviceImpl
}

func (a *adapterImpl) GetMonitor() public.SessionMonitor {
	return a.monitorImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}