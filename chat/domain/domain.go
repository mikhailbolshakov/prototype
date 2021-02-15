package domain

import (
	"context"
)

type Who uint

const (
	ADMIN Who = iota
	BOT
	USER
)

type From struct {
	Who        Who
	ChatUserId string
}

type UserStatus struct {
	ChatUserId string
	Status     string
}

type SetUserStatusRequest struct {
	ChatUserId string
	Status     string
	From       *From
}

type GetUsersStatusesRequest struct {
	ChatUserIds []string
}

type GetUsersStatusesResponse struct {
	Statuses []*UserStatus
}

type CreateUserRequest struct {
	Username string
	Password string
	Email    string
}

type CreateUserResponse struct {
	Id string
}

type CreateClientChannelRequest struct {
	ChatUserId  string
	TeamName    string
	DisplayName string
	Name        string
	Subscribers []string
}

type CreateClientChannelResponse struct {
	ChannelId string
}

type GetChannelsForUserAndMembersRequest struct {
	UserId        string
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
	ToChatUserId   string
	ChannelId      string
	Ephemeral      bool
	From           *From
	Attachments    []*PostAttachment
	PredefinedPost *PredefinedPost
}

type SubscribeUserRequest struct {
	ChatUserId string
	ChannelId  string
}

type AskBotRequest struct {
	Message string
	From    string
}

type AskBotResponse struct {
	Found  bool
	Answer string
}

// predefined posts codes
const (
	TP_CLIENT_NEW_REQUEST             = "client.new-request"
	TP_CLIENT_NEW_MED_REQUEST         = "client.new-med-request"
	TP_CLIENT_NEW_LAW_REQUEST         = "client.new-law-request"
	TP_CLIENT_REQUEST_ASSIGNED        = "client.request-assigned"
	TP_CONSULTANT_REQUEST_ASSIGNED    = "consultant.request-assigned"
	TP_CLIENT_NEW_EXPERT_CONSULTATION = "client.new-expert-consultation"
	TP_EXPERT_NEW_EXPERT_CONSULTATION = "expert.new-expert-consultation"
	TP_CLIENT_NO_CONSULTANT           = "client.no-consultant-available"
	TP_TASK_SOLVED                    = "client.task-solved"
	TP_CLIENT_FEEDBACK                = "client.feedback"
)

const (
	STATUS_OUT_OF_OFFICE = "ooo"
	STATUS_OFFLINE       = "offline"
	STATUS_AWAY          = "away"
	STATUS_DND           = "dnd"
	STATUS_ONLINE        = "online"
)

var UserStatusMap = map[string]struct{}{
	STATUS_OUT_OF_OFFICE: {},
	STATUS_OFFLINE:       {},
	STATUS_AWAY:          {},
	STATUS_DND:           {},
	STATUS_ONLINE:        {},
}

type LoginRequest struct {
	UserId     string
	Username   string
	ChatUserId string
}

type LoginResponse struct {
	ChatSessionId string
}

type LogoutRequest struct {
	ChatUserId string
}

type Service interface {
	GetUsersStatuses(ctx context.Context, rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateUser(ctx context.Context, rq *CreateUserRequest) (*CreateUserResponse, error)
	CreateClientChannel(ctx context.Context, rq *CreateClientChannelRequest) (*CreateClientChannelResponse, error)
	GetChannelsForUserAndMembers(ctx context.Context, rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	SubscribeUser(ctx context.Context, rq *SubscribeUserRequest) error
	DeleteUser(ctx context.Context, userId string) error
	AskBot(ctx context.Context, request *AskBotRequest) (*AskBotResponse, error)
	Posts(ctx context.Context, posts []*Post) error
	SetStatus(ctx context.Context, rq *SetUserStatusRequest) error
	Login(ctx context.Context, rq *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, rq *LogoutRequest) error
}
