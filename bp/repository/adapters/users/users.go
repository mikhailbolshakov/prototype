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

func (u *serviceImpl) Get(id string) *pb.User {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	user, _ := u.UsersClient.Get(ctx, &pb.GetByIdRequest{Id: id})
	return user
}

func (u *serviceImpl) Activate(userId string) (*pb.User, error) {
	return u.UsersClient.Activate(context.Background(), &pb.ActivateRequest{UserId: userId})
}

func (u *serviceImpl) Delete(userId string) (*pb.User, error) {
	return u.UsersClient.Delete(context.Background(), &pb.DeleteRequest{UserId: userId})
}

func (u *serviceImpl) SetClientDetails(userId string, details *pb.ClientDetails) (*pb.User, error) {
	return u.UsersClient.SetClientDetails(context.Background(), &pb.SetClientDetailsRequest{UserId: userId, ClientDetails: details})
}

func (u *serviceImpl) SetMMUserId(userId, mmId string) (*pb.User, error) {
	return u.UsersClient.SetMMUserId(context.Background(), &pb.SetMMIdRequest{UserId: userId, MMId: mmId})
}

func (u *serviceImpl) SetKKUserId(userId, kkId string) (*pb.User, error) {
	return u.UsersClient.SetKKUserId(context.Background(), &pb.SetKKIdRequest{UserId: userId, KKId: kkId})
}

func (u *serviceImpl) GetByMMId(mmUserId string) (*pb.User, error) {
	return u.UsersClient.GetByMMId(context.Background(), &pb.GetByMMIdRequest{MMId: mmUserId})
}