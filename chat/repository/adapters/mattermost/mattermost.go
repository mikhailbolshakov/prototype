package mattermost

import (
	"context"
	"fmt"
	"github.com/adacta-ru/mattermost-server/v6/model"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
)

type serviceImpl struct {
	hub ChatSessionHub
	cfg *kitConfig.Config
}

func newImpl(hub ChatSessionHub) *serviceImpl {
	m := &serviceImpl{
		hub: hub,
	}
	return m
}

func (s *serviceImpl) getClient(f *domain.From) (*Client, error) {

	var cl *Client
	switch f.Who {
	case domain.ADMIN:
		cl = s.hub.AdminSession().Client()
	case domain.BOT:
		cl = s.hub.BotSession().Client()
	case domain.USER:
		sess := s.hub.GetByChatUserId(f.ChatUserId)
		if sess == nil {
			return nil, fmt.Errorf("open session for chat user %s not found", f.ChatUserId)
		}
		cl = sess.Client()
	}

	if cl == nil {
		return nil, fmt.Errorf("mattermost connection not valid")
	}

	return cl, nil
}

// TODO: I don't like. We have to pass config on new somehow
func (s *serviceImpl) setConfig(cfg *kitConfig.Config) {
	s.cfg = cfg
}

func (s *serviceImpl) CreateUser(ctx context.Context, rq *domain.CreateUserRequest) (string, error) {

	if rq.Password == "" {
		rq.Password = s.cfg.Mattermost.DefaultPassword
	}

	userId, err := s.hub.AdminSession().Client().CreateUser(s.cfg.Mattermost.Team, rq.Username, rq.Password, rq.Email)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (s *serviceImpl) CreateClientChannel(ctx context.Context, rq *domain.CreateClientChannelRequest) (string, error) {

	if rq.TeamName == "" {
		rq.TeamName = s.cfg.Mattermost.Team
	}

	chId, err := s.hub.AdminSession().Client().CreateUserChannel("P", rq.TeamName, rq.ChatUserId, rq.DisplayName, rq.Name)
	if err != nil {
		return "", err
	}

	return chId, nil
}

func (s *serviceImpl) SubscribeUser(ctx context.Context, userId, channelId string) error {
	return s.hub.AdminSession().Client().SubscribeUser(channelId, userId)
}

func (s *serviceImpl) Post(ctx context.Context, post *domain.Post) error {

	cl, err := s.getClient(post.From)
	if err != nil {
		return err
	}

	if post.Attachments == nil {
		post.Attachments = []*domain.PostAttachment{}
	}
	att := s.convertAttachments(post.Attachments)

	if post.Ephemeral {
		return cl.CreateEphemeralPost(post.ChannelId, post.ToChatUserId, post.Message, att)
	} else {
		return cl.CreatePost(post.ChannelId, post.Message, att)
	}

}

func (s *serviceImpl) GetUserStatuses(ctx context.Context, rq *domain.GetUsersStatusesRequest) (*domain.GetUsersStatusesResponse, error) {

	rs := &domain.GetUsersStatusesResponse{
		Statuses: []*domain.UserStatus{},
	}

	if statuses, err := s.hub.AdminSession().Client().GetUsersStatuses(rq.ChatUserIds); err == nil {

		for _, s := range statuses {
			rs.Statuses = append(rs.Statuses, &domain.UserStatus{
				ChatUserId: s.UserId,
				Status: s.Status,
			})
		}

	} else {
		return rs, err
	}
	return rs, nil
}

func (s *serviceImpl) CreateDirectChannel(ctx context.Context, chatUserId1, chatUserId2 string) (string, error) {
	chId, err := s.hub.AdminSession().Client().CreateDirectChannel(chatUserId1, chatUserId2)
	if err != nil {
		return "", err
	}
	return chId, nil
}

func (s *serviceImpl) GetChannelsForUserAndMembers(ctx context.Context, rq *domain.GetChannelsForUserAndMembersRequest) ([]string, error) {
	return s.hub.AdminSession().Client().GetChannelsForUserAndMembers(rq.UserId, s.cfg.Mattermost.Team, rq.MemberUserIds)
}

func (s *serviceImpl) DeleteUser(ctx context.Context, chatUserId string) error {
	return s.hub.AdminSession().Client().DeleteUser(chatUserId)
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

func (s *serviceImpl) SetUserStatus(ctx context.Context, chatUserId, status string, from *domain.From) error {

	if cl, err := s.getClient(from); err != nil {
		return cl.UpdateStatus(chatUserId, status)
	} else {
		return err
	}

}

func (s *serviceImpl) Login(ctx context.Context, userId, username, chatUserId string) (string, error) {
	return s.hub.NewSession(ctx, userId, username, chatUserId)
}

func (s *serviceImpl) Logout(ctx context.Context, chatUserId string) error {
	return s.hub.Logout(ctx, chatUserId)
}
