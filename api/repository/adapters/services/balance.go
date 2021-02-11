package services

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
)



type balanceServiceImpl struct {
	pb.BalanceServiceClient
}

func newBalanceImpl() *balanceServiceImpl {
	a := &balanceServiceImpl{}
	return a
}

func (u *balanceServiceImpl) Add(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.Add(ctx, rq)
}

func (u *balanceServiceImpl) GetBalance(ctx context.Context, rq *pb.GetBalanceRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.GetBalance(ctx, rq)
}

func (u *balanceServiceImpl) WriteOff(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.WriteOff(ctx, rq)
}

func (u *balanceServiceImpl) Lock(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.Lock(ctx, rq)
}

func (u *balanceServiceImpl) CancelLock(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return u.BalanceServiceClient.CancelLock(ctx, rq)
}
