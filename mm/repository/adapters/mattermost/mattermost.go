package mattermost

import (
	"github.com/adacta-ru/mattermost-server/v6/model"
	"gitlab.medzdrav.ru/prototype/kit/chat/mattermost"
	"gitlab.medzdrav.ru/prototype/kit/queue"
)

type NewPostMessageHandler func(post *Post)

type Service interface {
	CreateUser(user *CreateUserRequest) (*CreateUserResponse, error)
	CreateClientChannel(rq *CreateClientChannelRequest) (*CreateChannelResponse, error)
	SubscribeUser(rq *SubscribeUserRequest) (*SubscribeUserResponse, error)
	CreateEphemeralPost(p *EphemeralPost) error
	CreatePost(p *Post) error
	GetUserStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateDirectChannel(rq *CreateDirectChannelRequest) (*CreateChannelResponse, error)
	// returns list of user's channels which have given active members (there might be more members and it's OK)
	GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	DeleteUser(userId string) error
}

type serviceImpl struct {
	client          *mattermost.Client
	queue           queue.Queue
	newPostsHandler NewPostMessageHandler
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

	userId, err := s.client.CreateUser(rq.TeamName, rq.Username, rq.Password, rq.Email)
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{userId}, nil
}

func (s *serviceImpl) CreateClientChannel(rq *CreateClientChannelRequest) (*CreateChannelResponse, error) {

	if rq.TeamName == "" {
		rq.TeamName = "rgs"
	}

	chId, err := s.client.CreateUserChannel("P", rq.TeamName, rq.ClientUserId, rq.DisplayName, rq.Name)
	if err != nil {
		return nil, err
	}

	return &CreateChannelResponse{chId}, nil
}

func (s *serviceImpl) SetNewPostMessageHandler(handler NewPostMessageHandler) {
	s.newPostsHandler = handler
}

func (s *serviceImpl) SubscribeUser(rq *SubscribeUserRequest) (*SubscribeUserResponse, error) {
	err := s.client.SubscribeUser(rq.ChannelId, rq.UserId)
	if err != nil {
		return nil, err
	}

	return &SubscribeUserResponse{}, nil
}

func (s *serviceImpl) CreateEphemeralPost(p *EphemeralPost) error {
	if p.Attachments == nil {
		p.Attachments = []*PostAttachment{}
	}
	return s.client.CreateEphemeralPost(p.ChannelId, p.UserId, p.Message, s.convertAttachments(p.Attachments))
}

func (s *serviceImpl) CreatePost(p *Post) error {
	if p.Attachments == nil {
		p.Attachments = []*PostAttachment{}
	}
	return s.client.CreatePost(p.ChannelId, p.Message, s.convertAttachments(p.Attachments))
}

func (s *serviceImpl) GetUserStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error) {

	rs := &GetUsersStatusesResponse{
		Statuses: []*UserStatus{},
	}

	if statuses, err := s.client.GetUsersStatuses(rq.UserIds); err == nil {

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
	chId, err := s.client.CreateDirectChannel(rq.UserId1, rq.UserId2)
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

	return s.client.GetChannelsForUserAndMembers(rq.UserId, teamName, rq.MemberUserIds)

}

func (s *serviceImpl) DeleteUser(userId string) error {
	return s.client.DeleteUser(userId)
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