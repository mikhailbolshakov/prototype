package domain

type UserStatus struct {
	UserId string
	Status string
}

type GetUsersStatusesRequest struct {
	UserIds []string
}

type GetUsersStatusesResponse struct {
	Statuses []*UserStatus
}

type CreateUserRequest struct {
	TeamName string
	Username string
	Password string
	Email    string
}

type CreateUserResponse struct {
	Id string
}

type CreateClientChannelRequest struct {
	ClientUserId string
	TeamName     string
	DisplayName  string
	Name         string
	Subscribers  []string
}

type CreateClientChannelResponse struct {
	ChannelId string
}

type GetChannelsForUserAndMembersRequest struct {
	UserId        string
	TeamName      string
	MemberUserIds []string
}

type SendTriggerPostRequest struct {
	TriggerPostCode string
	UserId          string
	ChannelId       string
	Params          map[string]interface{}
}

type SubscribeUserRequest struct {
	UserId    string
	ChannelId string
}