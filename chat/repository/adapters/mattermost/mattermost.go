package mattermost

import (
	"github.com/adacta-ru/mattermost-server/v6/model"
	"gitlab.medzdrav.ru/prototype/kit/chat/mattermost"
	"gitlab.medzdrav.ru/prototype/kit/queue"
)

// TODO: env
const (
	RGS_BOT_USERNAME     = "bot.rgs"
	RGS_BOT_ACCESS_TOKEN = "jg88x5sb63yk8ng6kcfkb37iho"
	MM_REST_URL          = "http://localhost:8065"
	MM_WS_URL            = "ws://localhost:8065"
	ADMIN_USERNAME       = "admin"
	ADMIN_PASSWORD       = "admin"
)

type NewPostMessageHandler func(post *Post)

type Service interface {
	CreateUser(user *CreateUserRequest) (*CreateUserResponse, error)
	CreateClientChannel(rq *CreateClientChannelRequest) (*CreateChannelResponse, error)
	SubscribeUser(rq *SubscribeUserRequest) (*SubscribeUserResponse, error)
	EphemeralPost(p *EphemeralPost) error
	Post(p *Post) error
	GetUserStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateDirectChannel(rq *CreateDirectChannelRequest) (*CreateChannelResponse, error)
	// returns list of user's channels which have given active members (there might be more members and it's OK)
	GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	DeleteUser(userId string) error
	PostFromBot(p *Post) error
	EphemeralPostFromBot(p *EphemeralPost) error
}

type serviceImpl struct {
	adminClient     *mattermost.Client
	queue           queue.Queue
	newPostsHandler NewPostMessageHandler
	botClient       *mattermost.Client
}

func newImpl(queue queue.Queue) *serviceImpl {
	m := &serviceImpl{
		queue: queue,
	}
	return m
}

func (s *serviceImpl) CreateUser(rq *CreateUserRequest) (*CreateUserResponse, error) {

	if rq.TeamName == "" {
		rq.TeamName = "rgs"
	}

	if rq.Password == "" {
		rq.Password = "12345"
	}

	userId, err := s.adminClient.CreateUser(rq.TeamName, rq.Username, rq.Password, rq.Email)
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{userId}, nil
}

func (s *serviceImpl) CreateClientChannel(rq *CreateClientChannelRequest) (*CreateChannelResponse, error) {

	if rq.TeamName == "" {
		rq.TeamName = "rgs"
	}

	chId, err := s.adminClient.CreateUserChannel("P", rq.TeamName, rq.ClientUserId, rq.DisplayName, rq.Name)
	if err != nil {
		return nil, err
	}

	return &CreateChannelResponse{chId}, nil
}

func (s *serviceImpl) SetNewPostMessageHandler(handler NewPostMessageHandler) {
	s.newPostsHandler = handler
}

func (s *serviceImpl) SubscribeUser(rq *SubscribeUserRequest) (*SubscribeUserResponse, error) {
	err := s.adminClient.SubscribeUser(rq.ChannelId, rq.UserId)
	if err != nil {
		return nil, err
	}

	return &SubscribeUserResponse{}, nil
}

func (s *serviceImpl) EphemeralPost(p *EphemeralPost) error {
	if p.Attachments == nil {
		p.Attachments = []*PostAttachment{}
	}
	return s.adminClient.CreateEphemeralPost(p.ChannelId, p.UserId, p.Message, s.convertAttachments(p.Attachments))
}

func (s *serviceImpl) Post(p *Post) error {
	if p.Attachments == nil {
		p.Attachments = []*PostAttachment{}
	}
	return s.adminClient.CreatePost(p.ChannelId, p.Message, s.convertAttachments(p.Attachments))
}

func (s *serviceImpl) GetUserStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error) {

	rs := &GetUsersStatusesResponse{
		Statuses: []*UserStatus{},
	}

	if statuses, err := s.adminClient.GetUsersStatuses(rq.UserIds); err == nil {

		for _, s := range statuses {
			rs.Statuses = append(rs.Statuses, &UserStatus{
				UserId: s.UserId,
				Status: s.Status,
			})
		}

	} else {
		return rs, err
	}
	return rs, nil
}

func (s *serviceImpl) CreateDirectChannel(rq *CreateDirectChannelRequest) (*CreateChannelResponse, error) {
	chId, err := s.adminClient.CreateDirectChannel(rq.UserId1, rq.UserId2)
	if err != nil {
		return nil, err
	}
	return &CreateChannelResponse{chId}, nil
}

func (s *serviceImpl) GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error) {

	teamName := "rgs"
	if rq.TeamName != "" {
		teamName = rq.TeamName
	}

	return s.adminClient.GetChannelsForUserAndMembers(rq.UserId, teamName, rq.MemberUserIds)

}

func (s *serviceImpl) DeleteUser(userId string) error {
	return s.adminClient.DeleteUser(userId)
}

func (s *serviceImpl) convertAttachments(attachments []*PostAttachment) []*model.SlackAttachment {

	var slackAttachments []*model.SlackAttachment

	for _, a := range attachments {

		sa := &model.SlackAttachment{
			Fallback:   a.Fallback,
			Color:      a.Color,
			Pretext:    a.Pretext,
			AuthorName: a.AuthorName,
			AuthorLink: a.AuthorLink,
			AuthorIcon: a.AuthorIcon,
			Title:      a.Title,
			TitleLink:  a.TitleLink,
			Text:       a.Text,
			ImageURL:   a.ImageURL,
			ThumbURL:   a.ThumbURL,
			Footer:     a.Footer,
			FooterIcon: a.FooterIcon,
		}

		if a.Fields != nil && len(a.Fields) > 0 {
			sa.Fields = []*model.SlackAttachmentField{}

			for _, f := range a.Fields {
				sa.Fields = append(sa.Fields, &model.SlackAttachmentField{
					Title: f.Title,
					Value: f.Value,
					Short: model.SlackCompatibleBool(f.Short),
				})

			}
		}

		slackAttachments = append(slackAttachments, sa)

	}

	return slackAttachments
}

func (s *serviceImpl) PostFromBot(p *Post) error {
	if p.Attachments == nil {
		p.Attachments = []*PostAttachment{}
	}
	return s.botClient.CreatePost(p.ChannelId, p.Message, s.convertAttachments(p.Attachments))
}

func (s *serviceImpl) EphemeralPostFromBot(p *EphemeralPost) error {
	if p.Attachments == nil {
		p.Attachments = []*PostAttachment{}
	}
	return s.botClient.CreateEphemeralPost(p.ChannelId, p.UserId, p.Message, s.convertAttachments(p.Attachments))
}
