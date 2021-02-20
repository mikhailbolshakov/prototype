package session

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.medzdrav.ru/prototype/api/public"
	"gitlab.medzdrav.ru/prototype/kit"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/proto"
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
	GetOutgoingWsEventsHandler() listener.QueueMessageHandler
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

	l := log.L().Cmp("hub").Mth("new-session")

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

	s := newSession(usr.Id, usr.Username, usr.MMId, chatSessionId, h.chatService, h.queue).setJWT(jwt)
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

	l.F(log.FF{"session-id": sessionId}).Inf("new session")
	l.TrcF("session details %v", s)

	return &NewSessionResponse{SessionId: sessionId}, nil
}

func (h *hubImpl) GetOutgoingWsEventsHandler() listener.QueueMessageHandler {

	return func(msg []byte) error {

		var wsEventMsgPl *proto.OutgoingWsEventQueueMessagePayload
		ctx, err := queue.Decode(context.Background(), msg, &wsEventMsgPl)
		if err != nil {
			return err
		}

		l := log.L().Pr("queue").Cmp("hub").Mth("ws-event-handler").C(ctx)
		l.TrcF("ws message details %s", string(msg))

		if wsEventMsgPl == nil {
			return fmt.Errorf("invalid message")
		}

		h.RLock()
		defer h.RUnlock()

		usrSessions := h.GetByUserId(wsEventMsgPl.UserId)
		if usrSessions == nil || len(usrSessions) == 0 {
			l.InfF("cannot forward to ws. no user session user=%s", wsEventMsgPl.UserId)
			return nil
		}

		// set message Id if not set
		if wsEventMsgPl.WsEvent.Id == "" {
			wsEventMsgPl.WsEvent.Id = kit.NewId()
		}

		// set correlationId as requestId from the ctx is not set
		if wsEventMsgPl.WsEvent.CorrelationId == "" {
			if r, err := kitContext.MustRequest(ctx); err == nil {
				wsEventMsgPl.WsEvent.CorrelationId = r.GetRequestId()
			} else {
				return err
			}
		}

		for _, s := range usrSessions {
			if err := s.sendWsMessage(wsEventMsgPl.WsEvent); err != nil {
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

	s := h.GetById(sessionId)

	if s != nil {

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