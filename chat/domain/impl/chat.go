package impl

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/chat/logger"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/log"
)

type serviceImpl struct {
	mmService   domain.MattermostService
}

func NewService(mmService domain.MattermostService) domain.Service {

	s := &serviceImpl{
		mmService: mmService,
	}

	return s
}

func (s *serviceImpl) l() log.CLogger {
	return logger.L().Cmp("chat")
}

func (s *serviceImpl) validateFrom(f *domain.From) error {

	switch f.Who {
	case domain.ADMIN, domain.BOT, domain.USER:
	default: return fmt.Errorf("invalid FROM value")
	}

	if f.Who == domain.USER && f.ChatUserId == "" {
		return fmt.Errorf("chat user Id must be specified for USER type")
	}

	return nil
}

func (s *serviceImpl) GetChannelsForUserAndMembers(ctx context.Context, rq *domain.GetChannelsForUserAndMembersRequest) ([]string, error) {

	l := s.l().C(ctx).Mth("get-channels-for-user").F(log.FF{"username": rq.UserId, "members": rq.MemberUserIds})

	r, err := s.mmService.GetChannelsForUserAndMembers(ctx, rq)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}
	l.DbgF("found %d items", len(r))

	return r, nil
}

func (s *serviceImpl) GetUsersStatuses(ctx context.Context, rq *domain.GetUsersStatusesRequest) (*domain.GetUsersStatusesResponse, error) {

	l := s.l().C(ctx).Mth("get-user-statuses").F(log.FF{"userIds": rq.ChatUserIds})

	r, err := s.mmService.GetUserStatuses(ctx, rq)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}
	l.DbgF("got %d items", len(r.Statuses))

	return r, nil
}

func (s *serviceImpl) CreateUser(ctx context.Context, rq *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {

	l := s.l().C(ctx).Mth("create-user").F(log.FF{"username": rq.Username})

	userId, err := s.mmService.CreateUser(ctx, rq)

	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	l.Dbg("created")

	return &domain.CreateUserResponse{Id: userId}, nil
}

func (s *serviceImpl) CreateClientChannel(ctx context.Context, rq *domain.CreateClientChannelRequest) (*domain.CreateClientChannelResponse, error) {

	l := s.l().C(ctx).Mth("create-channel").F(log.FF{"userId": rq.ChatUserId})

	channelId, err := s.mmService.CreateClientChannel(ctx, rq)
	if err != nil {
		l.E(err).St().Err("channel creation")
		return nil, err
	}
	l.F(log.FF{"channel": channelId}).Dbg("created")

	if rq.Subscribers != nil && len(rq.Subscribers) > 0 {
		for _, sbUserId := range rq.Subscribers {
			err = s.mmService.SubscribeUser(ctx, sbUserId, channelId)
			if err != nil {
				l.E(err).St().ErrF("subscription, userId=%s", sbUserId)
				return nil, err
			}
		}
	}
	l.DbgF("subscribed %d users", len(rq.Subscribers))

	return &domain.CreateClientChannelResponse{ChannelId: channelId}, nil
}

func (s *serviceImpl) SubscribeUser(ctx context.Context, rq *domain.SubscribeUserRequest) error {

	l := s.l().C(ctx).Mth("subscribe-user").F(log.FF{"userId": rq.ChatUserId, "channel": rq.ChannelId})

	err := s.mmService.SubscribeUser(ctx, rq.ChatUserId, rq.ChannelId)
	if err != nil {
		l.E(err).St().Err()
		return nil
	}

	l.Dbg("subscribed")
	return nil

}

func (s *serviceImpl) DeleteUser(ctx context.Context, userId string) error {

	l := s.l().C(ctx).Mth("delete-user").F(log.FF{"userId": userId})

	err := s.mmService.DeleteUser(ctx, userId)
	if err != nil {
		l.E(err).St().Err()
		return nil
	}

	l.Dbg("deleted")
	return nil

}

func (s *serviceImpl) Posts(ctx context.Context, posts []*domain.Post) error {

	l := s.l().C(ctx).Mth("posts")

	var err error
	for _, post := range posts {

		l.TrcF("%s", kit.MustJson(post))

		if post.Ephemeral && post.ToChatUserId == "" {
			err := fmt.Errorf("recipient user id must be specified for an ephemeral post")
			l.E(err).St().Err()
			return err
		}

		if err := s.validateFrom(post.From); err != nil {
			l.E(err).St().Err("validation")
			return err
		}

		if post.PredefinedPost != nil && post.PredefinedPost.Code != "" {
			post, err = s.predefinedPost(ctx, post)
			if err != nil {
				l.E(err).St().Err("predefined post")
				return err
			}
		}

		if err := s.mmService.Post(ctx, post); err != nil {
			l.E(err).St().Err("post")
			return err
		}

	}

	l.DbgF("posted %d posts", len(posts))

	return nil
}

func (s *serviceImpl) SetStatus(ctx context.Context, rq *domain.SetUserStatusRequest) error {

	l := s.l().C(ctx).Mth("set-status")

	if _, ok := domain.UserStatusMap[rq.Status]; !ok {
		err := fmt.Errorf("not valid status %s", rq.Status)
		l.E(err).St().Err()
		return err
	}

	if err := s.validateFrom(rq.From); err != nil {
		l.E(err).St().Err("validation")
		return err
	}

	err := s.mmService.SetUserStatus(ctx, rq.ChatUserId, rq.Status, rq.From)
	if err != nil {
		l.E(err).St().Err("set status")
		return err
	}

	return nil
}

func (s *serviceImpl) Login(ctx context.Context, rq *domain.LoginRequest) (*domain.LoginResponse, error) {

	l := s.l().C(ctx).Mth("login").F(log.FF{"username": rq.Username, "chatUser": rq.ChatUserId})

	sess, err := s.mmService.Login(ctx, rq.UserId, rq.Username, rq.ChatUserId)
	if err != nil {
		l.E(err).St().Err(err)
		return nil, err
	}
	l.DbgF("session %s", sess)
	return &domain.LoginResponse{ChatSessionId: sess}, nil
}

func (s *serviceImpl) Logout(ctx context.Context, rq *domain.LogoutRequest) error {
	l := s.l().C(ctx).Mth("logout").F(log.FF{"chatUser": rq.ChatUserId})
	err := s.mmService.Logout(ctx, rq.ChatUserId)
	if err != nil {
		l.E(err).St().Err()
		return err
	}
	l.Dbg("ok")
	return nil
}