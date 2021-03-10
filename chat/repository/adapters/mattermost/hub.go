package mattermost

import (
	"context"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/chat/logger"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"sync"
)

type NewChatSessionResponse struct {
	ChatSessionId string
}

type ChatSessionHub interface {
	Init(ctx context.Context) error
	NewSession(ctx context.Context, userId, username, chatUserId string) (string, error)
	Logout(ctx context.Context, chatUserId string) error
	GetById(id string) ChatSession
	GetByUserId(userId string) ChatSession
	GetByChatUserId(chatUserId string) ChatSession
	Close(ctx context.Context)
	AdminSession() ChatSession
	BotSession() ChatSession
}

type hubImpl struct {
	sync.RWMutex
	sessions                 map[string]ChatSession
	userSessions             map[string]ChatSession
	chatUserSessions         map[string]ChatSession
	adminSession, botSession ChatSession
	userService              domain.UserService
	cfg                      *config.Config
	queue                    queue.Queue
}

func NewHub(cfg *config.Config, userService domain.UserService, queue queue.Queue) ChatSessionHub {

	h := &hubImpl{
		userService:      userService,
		sessions:         map[string]ChatSession{},
		userSessions:     map[string]ChatSession{},
		chatUserSessions: map[string]ChatSession{},
		cfg:              cfg,
		queue:            queue,
	}

	return h
}

func (h *hubImpl) l() log.CLogger {
	return logger.L().Cmp("mm-hub")
}

func (h *hubImpl) Init(ctx context.Context) error {

	l := h.l().Mth("init")

	var err error

	h.adminSession, err = newAdminSession(ctx, h.cfg)
	if err != nil {
		return err
	}

	l.DbgF("admin session %s", h.adminSession.GetId())

	h.botSession, err = newBotSession(ctx, h.cfg)
	if err != nil {
		return err
	}

	l.DbgF("bot session %s", h.adminSession.GetId())

	return nil
}

func (h *hubImpl) NewSession(ctx context.Context, userId, username, chatUserId string) (string, error) {

	l := h.l().Mth("new-session").C(ctx).F(log.FF{"user": username})

	// we hold the only session for the userId
	h.RLock()
	s, ok := h.userSessions[userId]
	h.RUnlock()
	if ok {
		return s.GetId(), nil
	}

	s, err := newSession(ctx, userId, username, chatUserId, h.cfg, h.queue)
	if err != nil {
		return "", err
	}
	sessionId := s.GetId()

	func() {

		h.Lock()
		defer h.Unlock()

		h.sessions[sessionId] = s
		h.userSessions[userId] = s
		h.chatUserSessions[chatUserId] = s

	}()

	s.Open(ctx)

	l.F(log.FF{"sid": sessionId}).Dbg("logged in")

	return sessionId, nil
}

func (h *hubImpl) AdminSession() ChatSession {
	return h.adminSession
}

func (h *hubImpl) BotSession() ChatSession {
	return h.botSession
}

func (h *hubImpl) Logout(ctx context.Context, chatUserId string) error {

	l := h.l().Mth("logout").C(ctx).F(log.FF{"chat-user": chatUserId})

	h.Lock()
	defer h.Unlock()

	if s, ok := h.chatUserSessions[chatUserId]; ok {
		s.Close(ctx)
		delete(h.sessions, s.GetId())
		delete(h.userSessions, s.GetUserId())
		delete(h.chatUserSessions, s.GetChatUserId())
	} else {
		l.Warn("no session found")
		return nil
	}

	l.Dbg("logged out")

	return nil

}

func (h *hubImpl) GetById(id string) ChatSession {

	h.RLock()
	defer h.RUnlock()
	if s, ok := h.sessions[id]; ok {
		return s
	}

	return nil
}

func (h *hubImpl) GetByUserId(userId string) ChatSession {

	h.RLock()
	defer h.RUnlock()
	if s, ok := h.userSessions[userId]; ok {
		return s
	}

	return nil
}

func (h *hubImpl) GetByChatUserId(chatUserId string) ChatSession {
	h.RLock()
	defer h.RUnlock()
	if s, ok := h.chatUserSessions[chatUserId]; ok {
		return s
	}

	return nil
}

func (h *hubImpl) Close(ctx context.Context) {

	l := h.l().Mth("close")

	h.Lock()
	defer h.Unlock()

	h.adminSession.Close(ctx)
	h.botSession.Close(ctx)

	for _, s := range h.userSessions {
		s.Close(ctx)
	}

	h.sessions = nil
	h.userSessions = nil

	l.Dbg("closed")

}
