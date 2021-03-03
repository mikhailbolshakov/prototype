package domain

import (
	"context"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	userPb "gitlab.medzdrav.ru/prototype/proto/users"
)

type SessionStorage interface {}

type CfgService interface {
	Get(ctx context.Context) (*kit.Config, error)
}

type UserService interface {
	Get(ctx context.Context, id string) *userPb.User
}

type ChatService interface {
	SetStatus(ctx context.Context, userId, status string) error
	Login(ctx context.Context, userId, username, chatUserId string) (string, error)
	Logout(ctx context.Context, chatUserId string) error
}

type Metrics interface {
	SessionsInc()
	SessionsDec()
	ConnectedUsersInc()
	ConnectedUsersDec()
	SessionsCount() int
	ConnectedUsersCount() int
}