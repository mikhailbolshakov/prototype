package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto"
	"gitlab.medzdrav.ru/prototype/sessions/logger"
	"sync"
	"sync/atomic"
	"time"
)

type session interface {
	setWs(ws Ws) session
	sendWsMessage(msg *proto.WsMessage) error
	setJWT(jwt *auth.JWT) session
	isWs() bool
	getId() string
	getUserId() string
	getUsername() string
	getChatUserId() string
	getAccessToken() string
	getChatSessionId() string
	getStartAt() time.Time
	getSentWsMessages() uint32
	close(ctx context.Context)
}

type sessionImpl struct {
	sync.RWMutex
	id             string
	accessToken    string
	refreshToken   string
	expiresIn      int
	userId         string
	username       string
	chatUserId     string
	chatSessionId  string
	startAt        time.Time
	sentWsMessages uint32
	ws             Ws
	queue          queue.Queue
}

func newSession(userId, username, chatUserId string, chatSessionId string, queue queue.Queue) session {

	s := &sessionImpl{
		id:            kit.NewId(),
		userId:        userId,
		username:      username,
		chatUserId:    chatUserId,
		chatSessionId: chatSessionId,
		startAt:       time.Now().UTC(),
		queue:         queue,
	}

	return s
}

func (h *sessionImpl) l() log.CLogger {
	return logger.L().Cmp("session")
}

func (s *sessionImpl) getId() string {
	return s.id
}

func (s *sessionImpl) getUserId() string {
	return s.userId
}

func (s *sessionImpl) getChatUserId() string {
	return s.chatUserId
}

func (s *sessionImpl) getChatSessionId() string {
	// lock because chat session id isn't immutable
	// it could be changed if reconnect occurs
	s.RLock()
	defer s.RUnlock()
	return s.chatSessionId
}

func (s *sessionImpl) getStartAt() time.Time {
	return s.startAt
}

func (s *sessionImpl) getSentWsMessages() uint32 {
	return atomic.LoadUint32(&s.sentWsMessages)
}

func (s *sessionImpl) getUsername() string {
	return s.username
}

func (s *sessionImpl) getAccessToken() string {
	return s.accessToken
}

func (s *sessionImpl) forwardIncomingWsMessage(msg []byte) {

	l := s.l().Mth("fwd-incoming-msg").Dbg().TrcF("%s", string(msg))

	var wsMessage *proto.WsMessage
	err := json.Unmarshal(msg, &wsMessage)
	if err != nil {
		l.E(err).St().ErrF("invalid format")
		return
	}

	// build context
	rCtx := kitContext.NewRequestCtx().
		Ws().
		WithSessionId(s.getId()).
		WithUser(s.getUserId(), s.getUsername()).
		WithChatUserId(s.getChatUserId())

	if wsMessage.Id != "" {
		rCtx.WithRequestId(wsMessage.Id)
	} else {
		rCtx.WithNewRequestId()
	}

	ctx := rCtx.ToContext(context.Background())

	l.C(ctx)

	// define topic as template + message type
	topic := fmt.Sprintf(proto.QUEUE_TOPIC_INCOMING_WS_TEMPLATE, wsMessage.MessageType)

	qMsg := &queue.Message{
		Ctx:     rCtx,
		Payload: wsMessage,
	}

	// publish message
	if err := s.queue.Publish(ctx, queue.QUEUE_TYPE_AT_MOST_ONCE, topic, qMsg); err != nil {
		l.E(err).St().Err()
		return
	} else {
		l.DbgF("message forwarded, topic %s", topic)
	}

}

func (s *sessionImpl) wsListen(ws Ws) {

	closedEventChan := ws.wsClosedEvent()
	receivedEventChan := ws.receivedMessageEvent()
	for {
		select {
		case <-closedEventChan:
			s.Lock()
			s.ws = nil
			s.Unlock()
			return
		case msg := <-receivedEventChan:
			s.forwardIncomingWsMessage(msg)
		}
	}
}

func (s *sessionImpl) sendWsMessage(msg *proto.WsMessage) error {

	l := s.l().Mth("send-ws").Dbg()

	msgb, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// check if WS connection is open
	if s.isWs() {
		s.ws.send(msgb)
		// increment sent message counter
		atomic.AddUint32(&s.sentWsMessages, 1)
		l.F(log.FF{"user": s.userId}).TrcF("sent user=%s msg=%s", s.userId, string(msgb))
	} else {
		l.Warn("cannot send, no client connection")
	}
	return nil
}

func (s *sessionImpl) setWs(ws Ws) session {

	func() {
		s.Lock()
		defer s.Unlock()
		s.ws = ws
		ws.open()
	}()

	// start listening
	go s.wsListen(s.ws)

	return s
}

func (s *sessionImpl) setJWT(jwt *auth.JWT) session {
	s.accessToken = jwt.AccessToken
	s.refreshToken = jwt.RefreshToken
	s.expiresIn = jwt.ExpiresIn
	return s
}

func (s *sessionImpl) isWs() bool {
	s.RLock()
	defer s.RUnlock()
	return s.ws != nil
}

func (s *sessionImpl) close(ctx context.Context) {

	s.Lock()
	defer s.Unlock()

	if s.ws != nil {
		s.ws.close()
		s.ws = nil
	}

	s.l().Mth("session-close").C(ctx).Dbg("closed")

}
