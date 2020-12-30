package mm

type CreateUserRequest struct {
	TeamName string
	Username string
	Password string
	Email string
}

type CreateUserResponse struct {
	Id string
}

type CreateClientChannelRequest struct {
	ClientUserId string
	TeamName string
}

type CreateClientChannelResponse struct {
	ChannelId string
}

type Post struct {
	Id        string `json:"id"`
	CreateAt  int64  `json:"createAt"`
	UpdateAt  int64  `json:"updateAt"`
	EditAt    int64  `json:"editAt"`
	DeleteAt  int64  `json:"deleteAt"`
	UserId    string `json:"userId"`
	ChannelId string `json:"channelId"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}