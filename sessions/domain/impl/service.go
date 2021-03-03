package impl

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/sessions/domain"
	"golang.org/x/sync/errgroup"
)

type serviceImpl struct {
	hub         Hub
	userService domain.UserService
	chatService domain.ChatService
	auth        auth.Service
}

func NewSessionsService(h Hub, auth auth.Service, userService domain.UserService, chatService domain.ChatService) domain.SessionsService {
	return &serviceImpl{
		hub:         h,
		userService: userService,
		chatService: chatService,
		auth:        auth,
	}
}

func (s *serviceImpl) Login(ctx context.Context, rq *domain.LoginRequest) (*domain.LoginResponse, error) {

	l := log.L().C(ctx).Cmp("sessions").Mth("login").F(log.FF{"user": rq.Username, "chat": rq.ChatLogin}).Inf()

	usr := s.userService.Get(ctx, rq.Username)
	if usr == nil || usr.Id == "" {
		l.Err("no user found")
		return nil, fmt.Errorf("no user found %s", rq.Username)
	}

	var jwt *auth.JWT
	var chatSessionId string

	grp, _ := errgroup.WithContext(context.Background())
	grp.Go(func() error {
		var err error
		jwt, err = s.auth.AuthUser(&auth.User{
			UserName: rq.Username,
			Password: rq.Password,
		})
		return err
	})

	if rq.ChatLogin {
		grp.Go(func() error {
			var err error
			chatSessionId, err = s.chatService.Login(ctx, usr.Id, usr.Username, usr.MMId)
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	sid, err := s.hub.newSession(ctx, usr.Id, usr.Username, usr.MMId, chatSessionId, jwt)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	l.F(log.FF{"sid": sid}).Dbg()

	return &domain.LoginResponse{SessionId: sid}, nil
}

func (s *serviceImpl) Logout(ctx context.Context, rq *domain.LogoutRequest) error {

	l := log.L().C(ctx).Cmp("sessions").Mth("logout").F(log.FF{"user": rq.UserId}).Inf()

	sessions := s.hub.getByUserId(rq.UserId)
	if len(sessions) == 0 {
		l.Warn("no sessions found")
		return nil
	}

	for _, ss := range sessions {
		chatSid := ss.getChatSessionId()
		if chatSid != "" {
			if err := s.chatService.Logout(ctx, ss.getChatUserId()); err != nil {
				l.E(err).ErrF("chat session logout %s", chatSid)
			} else {
				l.DbgF("chat session %s logged out", chatSid)
			}
		}
	}

	err := s.hub.logout(ctx, rq.UserId)
	if err != nil {
		l.E(err).St().Err()
		return err
	}

	return nil
}

func (s *serviceImpl) toDomainSession(ss session) *domain.Session {
	if ss == nil {
		return nil
	}
	return &domain.Session{
		Id:            ss.getId(),
		UserId:        ss.getUserId(),
		Username:      ss.getUsername(),
		ChatUserId:    ss.getChatUserId(),
		ChatSessionId: ss.getChatSessionId(),
		LoginAt:       ss.getStartAt(),
	}
}

func (s *serviceImpl) Get(ctx context.Context, sid string) (*domain.Session, error) {
	ss := s.toDomainSession(s.hub.getById(sid))
	return ss, nil
}

func (s *serviceImpl) GetByUser(ctx context.Context, rq *domain.GetByUserRequest) ([]*domain.Session, error) {

	uid := rq.UserId

	if uid == "" && rq.Username != "" {
		usr := s.userService.Get(ctx, rq.Username)
		if usr == nil {
			return nil, fmt.Errorf("no user found %s", rq.Username)
		}
		uid = usr.Id
	}

	var rs []*domain.Session
	for _, ss := range s.hub.getByUserId(uid) {
		rs = append(rs, s.toDomainSession(ss))
	}

	return rs, nil
}

func (s *serviceImpl) AuthSession(ctx context.Context, sid string) (*domain.Session, error) {

	l := log.L().C(ctx).Cmp("sessions").Mth("auth").F(log.FF{"sid": sid}).Inf()

	ss := s.hub.getById(sid)
	if ss == nil {
		l.Err("no session found")
		return nil, fmt.Errorf("no session found")
	}

	err := s.auth.CheckAccessToken(ss.getAccessToken())
	if err != nil {
		return nil, err
	}

	return s.toDomainSession(ss), nil
}
