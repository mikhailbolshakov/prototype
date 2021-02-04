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

type CreateChannelResponse struct {
	ChannelId string
}

type PostAttachmentField struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
	Short bool        `json:"short"`
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
	// A unique Action ID. If not set, generated automatically.
	Id string `json:"id,omitempty"`

	// The type of the interactive element. Currently supported are
	// "select" and "button".
	Type string `json:"type,omitempty"`

	// The text on the button, or in the select placeholder.
	Name string `json:"name,omitempty"`

	// If the action is disabled.
	Disabled bool `json:"disabled,omitempty"`

	// Style defines a text and border style.
	// Supported values are "default", "primary", "success", "good", "warning", "danger"
	// and any hex color.
	Style string `json:"style,omitempty"`

	// DataSource indicates the data source for the select action. If left
	// empty, the select is populated from Options. Other supported values
	// are "users" and "channels".
	DataSource string `json:"data_source,omitempty"`

	// Options contains the values listed in a select dropdown on the post.
	Options []*PostActionOptions `json:"options,omitempty"`

	// DefaultOption contains the option, if any, that will appear as the
	// default selection in a select box. It has no effect when used with
	// other types of actions.
	DefaultOption string `json:"default_option,omitempty"`

	// Defines the interaction with the backend upon a user action.
	// Integration contains Context, which is private plugin data;
	// Integrations are stripped from Posts when they are sent to the
	// client, or are encrypted in a Cookie.
	Integration *PostActionIntegration `json:"integration,omitempty"`
	Cookie      string                 `json:"cookie,omitempty" db:"-"`
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
	Actions    []*PostAction          `json:"actions,omitempty"`
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
