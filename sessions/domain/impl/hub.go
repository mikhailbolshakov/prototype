package impl

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.medzdrav.ru/prototype/kit"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/proto"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"gitlab.medzdrav.ru/prototype/sessions/domain"
	"gitlab.medzdrav.ru/prototype/sessions/logger"
	"sync"
)

type hubImpl struct {
	http.Controller
	sync.RWMutex
	sessions     map[string]session
	userSessions map[string][]session
	httpServer   *http.Server
	cfg          *config.Config
	queue        queue.Queue
	metrics      domain.Metrics
}

type Hub interface {
	newSession(ctx context.Context, uid, username, chatUid, chatSid string, jwt *auth.JWT) (string, error)
	logout(ctx context.Context, userId string) error
	getById(id string) session
	getByUserId(userId string) []session
	setupWsConnection(sessionId string, wsConn *websocket.Conn) error
	GetOutgoingWsEventsHandler() listener.QueueMessageHandler
}

func NewHub(cfg *config.Config, srv *http.Server, metrics domain.Metrics) Hub {

	h := &hubImpl{
		httpServer:   srv,
		sessions:     map[string]session{},
		userSessions: map[string][]session{},
		cfg:          cfg,
		metrics:      metrics,
	}

	srv.SetWsUpgrader(newWsUpgrader(h))

	return h
}

func (h *hubImpl) l() log.CLogger {
	return logger.L().Cmp("hub")
}

func (h *hubImpl) newSession(ctx context.Context, uid, username, chatUid, chatSid string, jwt *auth.JWT) (string, error) {

	l := h.l().Mth("new-session").F(log.FF{"uid": uid, "un": username}).C(ctx).Dbg()

	s := newSession(uid, username, chatUid, chatSid, h.queue).setJWT(jwt)
	sessionId := s.getId()

	func() {

		h.Lock()
		defer h.Unlock()

		h.sessions[sessionId] = s
		h.metrics.SessionsInc()

		if us, ok := h.userSessions[uid]; ok {
			us = append(us, s)
			h.userSessions[uid] = us
		} else {
			h.userSessions[uid] = []session{s}
			h.metrics.ConnectedUsersInc()
		}

	}()

	l.F(log.FF{"sid": sessionId}).Inf().TrcF("%v", kit.MustJson(s))

	return sessionId, nil
}

func (h *hubImpl) GetOutgoingWsEventsHandler() listener.QueueMessageHandler {

	return func(msg []byte) error {

		var wsEventMsgPl *proto.OutgoingWsEventQueueMessagePayload
		ctx, err := queue.Decode(context.Background(), msg, &wsEventMsgPl)
		if err != nil {
			return err
		}

		l := h.l().Pr("queue").Mth("ws-event-handler").C(ctx)
		l.TrcF("ws message details %s", string(msg))

		if wsEventMsgPl == nil {
			return fmt.Errorf("invalid message")
		}

		h.RLock()
		defer h.RUnlock()

		usrSessions := h.getByUserId(wsEventMsgPl.UserId)
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

func (h *hubImpl) logout(ctx context.Context, userId string) error {

	l := h.l().Mth("logout").Dbg()

	h.Lock()
	defer h.Unlock()

	if _, ok := h.userSessions[userId]; !ok {
		err := fmt.Errorf("no sessions found for user %s", userId)
		l.E(err).Err()
		return err
	}

	for _, s := range h.userSessions[userId] {
		s.close(ctx)
		delete(h.sessions, s.getId())
		h.metrics.SessionsDec()
	}

	delete(h.userSessions, userId)
	h.metrics.ConnectedUsersDec()

	return nil

}

func (h *hubImpl) getById(id string) session {

	h.RLock()
	defer h.RUnlock()
	if s, ok := h.sessions[id]; ok {
		return s
	}

	return nil
}

func (h *hubImpl) getByUserId(userId string) []session {

	h.RLock()
	defer h.RUnlock()
	if s, ok := h.userSessions[userId]; ok {
		return s
	}

	return nil
}

func (h *hubImpl) setupWsConnection(sessionId string, wsConn *websocket.Conn) error {

	l := h.l().Mth("setup-ws-conn").F(log.FF{"sid": sessionId}).Dbg()

	s := h.getById(sessionId)

	if s != nil {

		if s.isWs() {
			err := fmt.Errorf("ws connection is already open for the session %s", sessionId)
			l.E(err).Err()
			return err
		} else {
			s.setWs(newWs(wsConn, s.getId(), s.getUserId()))
		}

	} else {
		err := fmt.Errorf("no active session found %s", sessionId)
		l.E(err).Err()
		return err
	}

	return nil
}

