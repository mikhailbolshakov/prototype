package session

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.medzdrav.ru/prototype/api/public"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"golang.org/x/sync/errgroup"
	net_http "net/http"
	"sync"
)

type NewSessionRequest struct {
	Username  string
	Password  string
	ChatLogin bool
}

type NewSessionResponse struct {
	SessionId string
}

type Hub interface {
	NewSession(context.Context, *NewSessionRequest) (*NewSessionResponse, error)
	Logout(ctx context.Context, userId string) error
	GetById(id string) Session
	GetByUserId(userId string) []Session
	SetupWsConnection(sessionId string, wsConn *websocket.Conn) error
	GetLoginRouteSetter() http.RouteSetter
	SessionMiddleware(next net_http.Handler) net_http.Handler
	NoSessionMiddleware(next net_http.Handler) net_http.Handler
	GetChatWsEventsHandler() listener.QueueMessageHandler
	GetMonitor() SessionMonitor
}

type hubImpl struct {
	http.Controller
	sync.RWMutex
	sessions     map[string]Session
	userSessions map[string][]Session
	auth         auth.Service
	userService  public.UserService
	chatService  public.ChatService
	httpServer   *http.Server
	cfg          *kitConfig.Config
	queue        queue.Queue
}

func NewHub(cfg *kitConfig.Config, srv *http.Server, auth auth.Service, userService public.UserService, chatService public.ChatService) Hub {

	h := &hubImpl{
		httpServer:   srv,
		auth:         auth,
		userService:  userService,
		chatService:  chatService,
		sessions:     map[string]Session{},
		userSessions: map[string][]Session{},
		cfg:          cfg,
	}

	srv.SetWsUpgrader(newWsUpgrader(h))

	return h
}

func (h *hubImpl) NewSession(ctx context.Context, rq *NewSessionRequest) (*NewSessionResponse, error) {

	usr := h.userService.Get(ctx, rq.Username)
	if usr == nil || usr.Id == "" {
		return nil, fmt.Errorf("no user found %s", rq.Username)
	}

	var jwt *auth.JWT
	var chatSessionId string

	grp, _ := errgroup.WithContext(context.Background())
	grp.Go(func() error {
		var err error
		jwt, err = h.auth.AuthUser(&auth.User{
			UserName: rq.Username,
			Password: rq.Password,
		})
		return err
	})

	if rq.ChatLogin {
		grp.Go(func() error {
			var err error
			chatSessionId, err = h.chatService.Login(ctx, usr.Id, usr.Username, usr.MMId)
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		return nil, err
	}

	s := newSession(usr.Id, usr.Username, usr.MMId, chatSessionId, h.chatService).setJWT(jwt)
	sessionId := s.getId()

	func() {

		h.Lock()
		defer h.Unlock()

		h.sessions[sessionId] = s

		if us, ok := h.userSessions[usr.Id]; ok {
			us = append(us, s)
			h.userSessions[usr.Id] = us
		} else {
			h.userSessions[usr.Id] = []Session{s}
		}

	}()

	log.TrcF("[hub] new session %v", s)

	return &NewSessionResponse{SessionId: sessionId}, nil
}

func (h *hubImpl) GetChatWsEventsHandler() listener.QueueMessageHandler {

	type ChatIncomingWsEventPayload struct {
		UserId     string
		ChatUserId string
		Event      string
		Data       map[string]interface{}
	}

	return func(msg []byte) error {

		var pl *ChatIncomingWsEventPayload
		_, err := queue.Decode(context.Background(), msg, &pl)
		if err != nil {
			return err
		}

		log.TrcF("[hub] chat ws message %v", string(msg))

		if pl == nil {
			return fmt.Errorf("[hub] invalid chat message")
		}

		h.RLock()
		defer h.RUnlock()

		usrSessions := h.GetByUserId(pl.UserId)
		if usrSessions == nil || len(usrSessions) == 0 {
			return fmt.Errorf("[hub] cannot send chat message to ws. user sessions not found. user=%s", pl.UserId)
		}

		msgToUser, _ := json.Marshal(pl)

		for _, s := range usrSessions {
			if err := s.sendWsMessage(msgToUser); err != nil {
				return err
			}
		}

		return nil
	}
}

func (h *hubImpl) Logout(ctx context.Context, userId string) error {
	h.Lock()
	defer h.Unlock()

	for _, s := range h.userSessions[userId] {
		s.close(ctx)
		delete(h.sessions, s.getId())
	}

	delete(h.userSessions, userId)

	return nil

}

func (h *hubImpl) GetById(id string) Session {

	h.RLock()
	defer h.RUnlock()
	if s, ok := h.sessions[id]; ok {
		return s
	}

	return nil
}

func (h *hubImpl) GetByUserId(userId string) []Session {

	h.RLock()
	defer h.RUnlock()
	if s, ok := h.userSessions[userId]; ok {
		return s
	}

	return nil
}

func (h *hubImpl) SetupWsConnection(sessionId string, wsConn *websocket.Conn) error {

	s, ok := func() (Session, bool) {
		h.RLock()
		defer h.RUnlock()
		s, ok := h.sessions[sessionId]
		return s, ok
	}()

	if ok {

		if s.isWs() {
			return fmt.Errorf("ws connection is already open for the session %s", sessionId)
		} else {
			s.setWs(newWs(wsConn, s.getId(), s.getUserId()))
		}

	} else {
		return fmt.Errorf("no active session found %s", sessionId)
	}

	return nil
}

func (h *hubImpl) GetMonitor() SessionMonitor {
	return h
}