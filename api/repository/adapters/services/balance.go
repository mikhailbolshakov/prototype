package services

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
)

type BalanceService interface {
	Add(rq *pb.ChangeServicesRequest) (*pb.UserBalance, error)
	GetBalance(rq *pb.GetBalanceRequest) (*pb.UserBalance, error)
	WriteOff(rq *pb.ChangeServicesRequest) (*pb.UserBalance, error)
	Lock(rq *pb.ChangeServicesRequest) (*pb.UserBalance, error)
	CancelLock(rq *pb.ChangeServicesRequest) (*pb.UserBalance, error)
}

type balanceServiceImpl struct {
	pb.BalanceServiceClient
}

func newBalanceImpl() *balanceServiceImpl {
	a := &balanceServiceImpl{}
	return a
}

func (u *balanceServiceImpl) Add(rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.Add(context.Background(), rq)
}

func (u *balanceServiceImpl) GetBalance(rq *pb.GetBalanceRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.GetBalance(context.Background(), rq)
}

func (u *balanceServiceImpl) WriteOff(rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.WriteOff(context.Background(), rq)
}

func (u *balanceServiceImpl) Lock(rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.Lock(context.Background(), rq)
}

func (u *balanceServiceImpl) CancelLock(rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.CancelLock(context.Background(), rq)
}
