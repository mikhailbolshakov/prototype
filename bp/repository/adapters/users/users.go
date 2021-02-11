package users

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)


type serviceImpl struct {
	pb.UsersClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) Get(ctx context.Context, id string) *pb.User {
	user, _ := u.UsersClient.Get(ctx, &pb.GetByIdRequest{Id: id})
	return user
}

func (u *serviceImpl) Activate(ctx context.Context, userId string) (*pb.User, error) {
	return u.UsersClient.Activate(ctx, &pb.ActivateRequest{UserId: userId})
}

func (u *serviceImpl) Delete(ctx context.Context, userId string) (*pb.User, error) {
	return u.UsersClient.Delete(ctx, &pb.DeleteRequest{UserId: userId})
}

func (u *serviceImpl) SetClientDetails(ctx context.Context, userId string, details *pb.ClientDetails) (*pb.User, error) {
	return u.UsersClient.SetClientDetails(ctx, &pb.SetClientDetailsRequest{UserId: userId, ClientDetails: details})
}

func (u *serviceImpl) SetMMUserId(ctx context.Context, userId, mmId string) (*pb.User, error) {
	return u.UsersClient.SetMMUserId(ctx, &pb.SetMMIdRequest{UserId: userId, MMId: mmId})
}

func (u *serviceImpl) SetKKUserId(ctx context.Context, userId, kkId string) (*pb.User, error) {
	return u.UsersClient.SetKKUserId(ctx, &pb.SetKKIdRequest{UserId: userId, KKId: kkId})
}

func (u *serviceImpl) GetByMMId(ctx context.Context, mmUserId string) (*pb.User, error) {
	return u.UsersClient.GetByMMId(ctx, &pb.GetByMMIdRequest{MMId: mmUserId})
}