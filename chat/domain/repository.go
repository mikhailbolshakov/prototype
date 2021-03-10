package domain

import (
	"context"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type ConfigService interface {
	Get() (*config.Config, error)
}

type MattermostService interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) (string, error)
	CreateClientChannel(ctx context.Context, rq *CreateClientChannelRequest) (string, error)
	SubscribeUser(ctx context.Context, userId, channelId string) error
	Post(ctx context.Context, post *Post) error
	GetUserStatuses(ctx context.Context, rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateDirectChannel(ctx context.Context, userId1, userId2 string) (string, error)
	// returns list of user's channels which have given active members (there might be more members and it's OK)
	GetChannelsForUserAndMembers(ctx context.Context, rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	DeleteUser(ctx context.Context, userId string) error
	SetUserStatus(ctx context.Context, chatUserId, status string, from *From) error
	Login(ctx context.Context, userId, username, chatUserId string) (string, error)
	Logout(ctx context.Context, chatUserId string) error
}

type UserService interface {
	Get(ctx context.Context, id string) *pb.User
}