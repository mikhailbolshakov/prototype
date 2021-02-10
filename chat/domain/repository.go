package domain

import kit "gitlab.medzdrav.ru/prototype/kit/config"

type ConfigService interface {
	Get() (*kit.Config, error)
}

type MattermostService interface {
	CreateUser(user *CreateUserRequest) (string, error)
	CreateClientChannel(rq *CreateClientChannelRequest) (string, error)
	SubscribeUser(userId, channelId string) error
	Post(channelId, message, toUserId string, ephemeral, fromBot bool, attachments []*PostAttachment) error
	GetUserStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateDirectChannel(userId1, userId2 string) (string, error)
	// returns list of user's channels which have given active members (there might be more members and it's OK)
	GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	DeleteUser(userId string) error
}
