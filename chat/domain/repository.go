package domain

import (
	"context"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
)

type ConfigService interface {
	Get() (*kit.Config, error)
}

type MattermostService interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) (string, error)
	CreateClientChannel(ctx context.Context, rq *CreateClientChannelRequest) (string, error)
	SubscribeUser(ctx context.Context, userId, channelId string) error
	Post(ctx context.Context, channelId, message, toUserId string, ephemeral, fromBot bool, attachments []*PostAttachment) error
	GetUserStatuses(ctx context.Context, rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateDirectChannel(ctx context.Context, userId1, userId2 string) (string, error)
	// returns list of user's channels which have given active members (there might be more members and it's OK)
	GetChannelsForUserAndMembers(ctx context.Context, rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	DeleteUser(ctx context.Context, userId string) error
}
