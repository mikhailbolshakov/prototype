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

type PredefinedPost struct {
	Code   string
	Params map[string]interface{}
}

type PostActionOptions struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type PostActionIntegration struct {
	URL     string                 `json:"url,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

type PostAction struct {
	Id            string                 `json:"id,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Disabled      bool                   `json:"disabled,omitempty"`
	Style         string                 `json:"style,omitempty"`
	DataSource    string                 `json:"data_source,omitempty"`
	Options       []*PostActionOptions   `json:"options,omitempty"`
	DefaultOption string                 `json:"default_option,omitempty"`
	Integration   *PostActionIntegration `json:"integration,omitempty"`
	Cookie        string                 `json:"cookie,omitempty" db:"-"`
}

type PostAttachment struct {
	Fallback   string
	Color      string
	Pretext    string
	AuthorName string
	AuthorLink string
	AuthorIcon string
	Title      string
	TitleLink  string
	Text       string
	ImageURL   string
	ThumbURL   string
	Footer     string
	FooterIcon string
	Actions    []*PostAction
}

type Post struct {
	Message        string
	ToUserId       string
	ChannelId      string
	Ephemeral      bool
	FromBot        bool
	Attachments    []*PostAttachment
	PredefinedPost *PredefinedPost
}

type SubscribeUserRequest struct {
	UserId    string
	ChannelId string
}

type AskBotRequest struct {
	Message string
	From    string
}

type AskBotResponse struct {
	Found  bool
	Answer string
}
