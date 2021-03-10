package domain

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"gitlab.medzdrav.ru/prototype/proto/config"
)

type UserStorage interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	GetByUsername(ctx context.Context, username string) *User
	GetByMMId(ctx context.Context, mmId string) *User
	Get(ctx context.Context, id string) *User
	Search(ctx context.Context, cr *SearchCriteria) (*SearchResponse, error)
	UpdateStatus(ctx context.Context, userId, status string, isDeleted bool) (*User, error)
	UpdateDetails(ctx context.Context, userId string, details string) (*User, error)
	UpdateMMId(ctx context.Context, userId, mmId string) (*User, error)
	UpdateKKId(ctx context.Context, userId, kkId string) (*User, error)
}

type ChatService interface {
	CreateUser(ctx context.Context, rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	CreateClientChannel(ctx context.Context, rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error)
	GetUsersStatuses(ctx context.Context, rq *pb.GetUsersStatusesRequest) (*pb.GetUserStatusesResponse, error)
}

type ConfigService interface {
	Get(ctx context.Context) (*config.Config, error)
}