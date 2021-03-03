package chat

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type serviceImpl struct {
	pb.UsersClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}


func (u *serviceImpl) SetStatus(ctx context.Context, chatUserId, status string) error {

	_, err := u.UsersClient.SetStatus(ctx, &pb.SetStatusRequest{
		From:       &pb.From{ Who: pb.From_ADMIN },
		UserStatus: &pb.UserStatus{
			Status:     status,
			ChatUserId: chatUserId,
		},
	})

	return err
}

func (u *serviceImpl) Login(ctx context.Context, userId, username, chatUserId string) (string, error) {
	lr, err := u.UsersClient.Login(ctx, &pb.LoginRequest{
		UserId:     userId,
		ChatUserId: chatUserId,
		Username:   username,
	})
	if err != nil {
		return "", err
	}
	return lr.ChatSessionId, nil
}

func (u *serviceImpl) Logout(ctx context.Context, chatUserId string) error {
	_, err := u.UsersClient.Logout(ctx, &pb.LogoutRequest{
		ChatUserId: chatUserId,
	})
	return err
}

