package mattermost

type SubscribeUserRequest struct {
	UserId    string
	ChannelId string
}

type SubscribeUserResponse struct{}

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
}

type CreateDirectChannelRequest struct {
	UserId1 string
	UserId2 string
}

type CreateChannelResponse struct {
	ChannelId string
}

type PostAttachmentField struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
	Short bool        `json:"short"`
}

type PostAttachment struct {
	Fallback   string                 `json:"fallback"`
	Color      string                 `json:"color"`
	Pretext    string                 `json:"pretext"`
	AuthorName string                 `json:"author_name"`
	AuthorLink string                 `json:"author_link"`
	AuthorIcon string                 `json:"author_icon"`
	Title      string                 `json:"title"`
	TitleLink  string                 `json:"title_link"`
	Text       string                 `json:"text"`
	Fields     []*PostAttachmentField `json:"fields"`
	ImageURL   string                 `json:"image_url"`
	ThumbURL   string                 `json:"thumb_url"`
	Footer     string                 `json:"footer"`
	FooterIcon string                 `json:"footer_icon"`
}

type Post struct {
	Id          string            `json:"id"`
	CreateAt    int64             `json:"createAt"`
	UpdateAt    int64             `json:"updateAt"`
	EditAt      int64             `json:"editAt"`
	DeleteAt    int64             `json:"deleteAt"`
	UserId      string            `json:"userId"`
	ChannelId   string            `json:"channelId"`
	Message     string            `json:"message"`
	Type        string            `json:"type"`
	Attachments []*PostAttachment `json:"attachments"`
}

type EphemeralPost struct {
	*Post
	UserId string
}

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

type GetChannelsForUserAndMembersRequest struct {
	UserId        string
	TeamName      string
	MemberUserIds []string
}
