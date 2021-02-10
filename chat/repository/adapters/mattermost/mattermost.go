package mattermost

import (
	"github.com/adacta-ru/mattermost-server/v6/model"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/kit/chat/mattermost"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
)

type serviceImpl struct {
	adminClient     *mattermost.Client
	botClient       *mattermost.Client
	cfg             *kitConfig.Config
}

func newImpl() *serviceImpl {
	m := &serviceImpl{}
	return m
}

// TODO: I don't like. We have to pass config on new somehow
func (s *serviceImpl) setConfig(cfg *kitConfig.Config) {
	s.cfg = cfg
}

func (s *serviceImpl) CreateUser(rq *domain.CreateUserRequest) (string, error) {

	if rq.TeamName == "" {
		rq.TeamName = s.cfg.Mattermost.Team
	}

	if rq.Password == "" {
		rq.Password = s.cfg.Mattermost.DefaultPassword
	}

	userId, err := s.adminClient.CreateUser(rq.TeamName, rq.Username, rq.Password, rq.Email)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (s *serviceImpl) CreateClientChannel(rq *domain.CreateClientChannelRequest) (string, error) {

	if rq.TeamName == "" {
		rq.TeamName = s.cfg.Mattermost.Team
	}

	chId, err := s.adminClient.CreateUserChannel("P", rq.TeamName, rq.ClientUserId, rq.DisplayName, rq.Name)
	if err != nil {
		return "", err
	}

	return chId, nil
}

func (s *serviceImpl) SubscribeUser(userId, channelId string) error {
	return s.adminClient.SubscribeUser(channelId, userId)
}

func (s *serviceImpl) Post(channelId, message, toUserId string, ephemeral, fromBot bool, attachments []*domain.PostAttachment) error {

	var client *mattermost.Client
	if fromBot {
		client = s.botClient
	} else {
		client = s.adminClient
	}

	if attachments == nil {
		attachments = []*domain.PostAttachment{}
	}
	att := s.convertAttachments(attachments)

	if ephemeral {
		return client.CreateEphemeralPost(channelId, toUserId, message, att)
	} else {
		return client.CreatePost(channelId, message, att)
	}

}

func (s *serviceImpl) GetUserStatuses(rq *domain.GetUsersStatusesRequest) (*domain.GetUsersStatusesResponse, error) {

	rs := &domain.GetUsersStatusesResponse{
		Statuses: []*domain.UserStatus{},
	}

	if statuses, err := s.adminClient.GetUsersStatuses(rq.UserIds); err == nil {

		for _, s := range statuses {
			rs.Statuses = append(rs.Statuses, &domain.UserStatus{
				UserId: s.UserId,
				Status: s.Status,
			})
		}

	} else {
		return rs, err
	}
	return rs, nil
}

func (s *serviceImpl) CreateDirectChannel(userId1, userId2 string) (string, error) {
	chId, err := s.adminClient.CreateDirectChannel(userId1, userId2)
	if err != nil {
		return "", err
	}
	return chId, nil
}

func (s *serviceImpl) GetChannelsForUserAndMembers(rq *domain.GetChannelsForUserAndMembersRequest) ([]string, error) {

	teamName := s.cfg.Mattermost.Team
	if rq.TeamName != "" {
		teamName = rq.TeamName
	}

	return s.adminClient.GetChannelsForUserAndMembers(rq.UserId, teamName, rq.MemberUserIds)

}

func (s *serviceImpl) DeleteUser(userId string) error {
	return s.adminClient.DeleteUser(userId)
}

func (s *serviceImpl) convertAttachments(attachments []*domain.PostAttachment) []*model.SlackAttachment {

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

		//if a.Fields != nil && len(a.Fields) > 0 {
		//	sa.Fields = []*model.SlackAttachmentField{}
		//
		//	for _, f := range a.Fields {
		//		sa.Fields = append(sa.Fields, &model.SlackAttachmentField{
		//			Title: f.Title,
		//			Value: f.Value,
		//			Short: model.SlackCompatibleBool(f.Short),
		//		})
		//
		//	}
		//}

		if a.Actions != nil && len(a.Actions) > 0 {
			sa.Actions = []*model.PostAction{}
			for _, act := range a.Actions {
				sAct := &model.PostAction{
					Id:            act.Id,
					Type:          act.Type,
					Name:          act.Name,
					Disabled:      act.Disabled,
					Style:         act.Style,
					DataSource:    act.DataSource,
					Options:       []*model.PostActionOptions{},
					DefaultOption: act.DefaultOption,
					Integration:   &model.PostActionIntegration{},
					Cookie:        act.Cookie,
				}

				if act.Integration != nil {
					sAct.Integration.URL = act.Integration.URL
					sAct.Integration.Context = act.Integration.Context
				}

				if act.Options != nil && len(act.Options) > 0 {
					for _, o := range act.Options {
						sAct.Options = append(sAct.Options, &model.PostActionOptions{
							Text:  o.Text,
							Value: o.Value,
						})
					}
				}

				sa.Actions = append(sa.Actions, sAct)
			}

		}

		slackAttachments = append(slackAttachments, sa)

	}

	return slackAttachments
}
