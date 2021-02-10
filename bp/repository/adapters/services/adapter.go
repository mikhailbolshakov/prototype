package services

import (
	"gitlab.medzdrav.ru/prototype/bp/domain"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
	GetBalanceService() domain.BalanceService
	GetDeliveryService() domain.DeliveryService
	Close()
}

type adapterImpl struct {
	balanceServiceImpl  *balanceServiceImpl
	deliveryServiceImpl *deliveryServiceImpl
	client              *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		balanceServiceImpl:  newBalanceImpl(),
		deliveryServiceImpl: newDeliveryImpl(),
	}
	return a
}

func (a *adapterImpl) Init(c *kitConfig.Config) error {
	cfg := c.Services["services"]
	cl, err := kitGrpc.NewClient(cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.balanceServiceImpl.BalanceServiceClient = pb.NewBalanceServiceClient(cl.Conn)
	a.deliveryServiceImpl.DeliveryServiceClient = pb.NewDeliveryServiceClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetBalanceService() domain.BalanceService {
	return a.balanceServiceImpl
}

func (a *adapterImpl) GetDeliveryService() domain.DeliveryService {
	return a.deliveryServiceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}