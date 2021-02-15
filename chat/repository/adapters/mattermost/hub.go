package mattermost

import (
	"context"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"sync"
)

type NewChatSessionResponse struct {
	ChatSessionId string
}

const INCOMING_WS_QUEUE_TOPIC = "mm.ws.event"

// TODO: move to config
// for the whole list of events refer to https://api.mattermost.com/#tag/WebSocket
var skippedWsEventTypes = map[string]struct{}{
	"hello": {},
	"config_changed": {},
	"delete_team": {},
	"leave_team": {},
	"license_changed": {},
	"plugin_disabled": {},
	"plugin_enabled": {},
	"plugin_statuses_changed": {},
	"preference_changed": {},
	"preferences_changed": {},
	"preferences_deleted": {},
	"response": {},
	"update_team": {},
	"user_added": {},
	"user_removed": {},
	"user_role_updated": {},
	"user_updated": {},
}

type ChatIncomingWsEventPayload struct {
	UserId     string
	ChatUserId string
	Event      string
	Data       map[string]interface{}
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
	cfg                      *kitConfig.Config
	queue                    queue.Queue
}

func NewHub(cfg *kitConfig.Config, userService domain.UserService, queue queue.Queue) ChatSessionHub {

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

func (h *hubImpl) Init(ctx context.Context) error {

	var err error

	h.adminSession, err = newAdminSession(ctx, h.cfg)
	if err != nil {
		return err
	}
	h.botSession, err = newBotSession(ctx, h.cfg)
	if err != nil {
		return err
	}

	return nil
}

func (h *hubImpl) NewSession(ctx context.Context, userId, username, chatUserId string) (string, error) {

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

	log.DbgF("[chat] user %s logged in session=%s", userId, sessionId)

	s.Open(ctx)

	return sessionId, nil
}

func (h *hubImpl) AdminSession() ChatSession {
	return h.adminSession
}

func (h *hubImpl) BotSession() ChatSession {
	return h.botSession
}

func (h *hubImpl) Logout(ctx context.Context, chatUserId string) error {
	h.Lock()
	defer h.Unlock()

	if s, ok := h.chatUserSessions[chatUserId]; ok {
		s.Close(ctx)
		delete(h.sessions, s.GetId())
		delete(h.userSessions, s.GetUserId())
		delete(h.chatUserSessions, s.GetChatUserId())
	} else {
		log.DbgF("[chat][logout] no session found for user %s", chatUserId)
		return nil
	}

	log.DbgF("[chat] user %s logged out", chatUserId)

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

	h.Lock()
	defer h.Unlock()

	h.adminSession.Close(ctx)
	h.botSession.Close(ctx)

	for _, s := range h.userSessions {
		s.Close(ctx)
	}

	h.sessions = nil
	h.userSessions = nil

	log.Dbg("[chat] hub closed")

}
