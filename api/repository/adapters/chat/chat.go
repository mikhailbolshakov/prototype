package chat

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/api/public"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type serviceImpl struct {
	pb.UsersClient
	pb.PostsClient
	userService public.UserService
}

func newImpl(userService public.UserService) *serviceImpl {
	a := &serviceImpl{
		userService: userService,
	}
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

func (u *serviceImpl) Post(ctx context.Context, fromUserId, channelId, message string) error {

	r, err := kitContext.MustRequest(ctx)
	if err != nil {
		return err
	}

	var chatUserId string
	if r.GetUserId() == fromUserId {
		chatUserId = r.GetChatUserId()
	} else {
		usr := u.userService.Get(ctx, fromUserId)
		if usr == nil {
			return fmt.Errorf("user not found %s", fromUserId)
		}
		chatUserId = usr.MMId
	}

	_, err = u.PostsClient.Post(ctx, &pb.PostRequest{Posts: []*pb.Post{{
		From:           &pb.From{Who: pb.From_USER, ChatUserId: chatUserId},
		Message:        message,
		ChannelId:      channelId,
	}}})
	return err
}

func (u *serviceImpl) EphemeralPost(ctx context.Context, fromUserId, toUserId, channelId, message string) error {

	r, err := kitContext.MustRequest(ctx)
	if err != nil {
		return err
	}

	var chatUserId string
	if r.GetUserId() == fromUserId {
		chatUserId = r.GetChatUserId()
	} else {
		usr := u.userService.Get(ctx, fromUserId)
		if usr == nil {
			return fmt.Errorf("user not found %s", fromUserId)
		}
		chatUserId = usr.MMId
	}

	toUsr := u.userService.Get(ctx, toUserId)
	if toUsr == nil {
		return fmt.Errorf("user not found %s", toUserId)
	}

	_, err = u.PostsClient.Post(ctx, &pb.PostRequest{Posts: []*pb.Post{{
		From:         &pb.From{Who: pb.From_USER, ChatUserId: chatUserId},
		Message:      message,
		ChannelId:    channelId,
		Ephemeral:    true,
		ToChatUserId: toUsr.MMId,
	}}})

	return err
}