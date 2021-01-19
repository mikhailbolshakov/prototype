package services

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
)

type Adapter interface {
	Init() error
	GetBalanceService() BalanceService
	GetDeliveryService() DeliveryService
}

type adapterImpl struct {
	balanceServiceImpl  *balanceServiceImpl
	deliveryServiceImpl *deliveryServiceImpl
	initialized         bool
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		balanceServiceImpl:  newBalanceImpl(),
		deliveryServiceImpl: newDeliveryImpl(),
		initialized:         false,
	}
	return a
}

func (a *adapterImpl) Init() error {
	cl, err := kitGrpc.NewClient("localhost", "50054")
	if err != nil {
		return err
	}
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
