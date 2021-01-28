package session

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/api/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/kit/chat/mattermost"
	"gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
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
	NewSession(*NewSessionRequest) (*NewSessionResponse, error)
	Logout(userId string) error
	GetById(id string) Session
	GetByUserId(userId string) Session
	SetWs(sessionId string, ws Ws) error
	GetLoginRouteSetter() http.RouteSetter
	SessionMiddleware(next net_http.Handler) net_http.Handler
}

type hubImpl struct {
	http.Controller
	sync.RWMutex
	sessions     map[string]Session
	userSessions map[string][]Session
	auth         auth.Service
	userService  users.Service
	httpServer   *http.Server
}

func NewHub(srv *http.Server, auth auth.Service, userService users.Service) Hub {

	h := &hubImpl{
		httpServer:   srv,
		auth:         auth,
		userService:  userService,
		sessions:     map[string]Session{},
		userSessions: map[string][]Session{},
	}

	srv.SetWsUpgrader(newWsUpgrader(h))

	return h
}

func (h *hubImpl) NewSession(rq *NewSessionRequest) (*NewSessionResponse, error) {

	usr := h.userService.Get(rq.Username)
	if usr == nil || usr.Id == "" {
		return nil, fmt.Errorf("no user found %s", rq.Username)
	}

	var jwt *auth.JWT
	var mmClient *mattermost.Client

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
			mmClient, err = mattermost.Login(&mattermost.Params{
				// TODO: env
				Url:     "http://localhost:8065",
				WsUrl:   "ws://localhost:8065",
				Account: rq.Username,
				Pwd:     rq.Password,
				OpenWS:  true,
			})
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		return nil, err
	}

	s := newSession(usr.Id, mmClient).setJWT(jwt)
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

	return &NewSessionResponse{SessionId: sessionId}, nil
}

func (h *hubImpl) Logout(userId string) error {
	h.Lock()
	defer h.Unlock()

	for _, s := range h.userSessions[userId] {
		s.close()
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

func (h *hubImpl) GetByUserId(userId string) Session {
	return nil
}

func (h *hubImpl) SetWs(sessionId string, ws Ws) error {

	s, ok := func() (Session, bool) {
		h.RLock()
		defer h.RUnlock()
		s, ok := h.sessions[sessionId]
		return s, ok
	}()

	if ok {
		s.setWs(ws)
	} else {
		return fmt.Errorf("no active session found %s", sessionId)
	}

	return nil
}
