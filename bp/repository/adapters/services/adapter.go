package services

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
)

type Adapter interface {
	Init() error
	GetBalanceService() BalanceService
	GetDeliveryService() DeliveryService
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

func (a *adapterImpl) Init() error {
	cl, err := kitGrpc.NewClient("localhost", "50054")
	if err != nil {
		return err
	}
	a.client = cl
	a.balanceServiceImpl.BalanceServiceClient = pb.NewBalanceServiceClient(cl.Conn)
	a.deliveryServiceImpl.DeliveryServiceClient = pb.NewDeliveryServiceClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetBalanceService() BalanceService {
	return a.balanceServiceImpl
}

func (a *adapterImpl) GetDeliveryService() DeliveryService {
	return a.deliveryServiceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}