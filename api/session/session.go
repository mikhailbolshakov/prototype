package session

import (
	"context"
	"gitlab.medzdrav.ru/prototype/api/public"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"sync"
	"sync/atomic"
	"time"
)

type Session interface {
	setWs(ws Ws) Session
	sendWsMessage(msg []byte) error
	setJWT(jwt *auth.JWT) Session
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
	id             string
	accessToken    string
	refreshToken   string
	expiresIn      int
	userId         string
	username       string
	chatUserId     string
	chatSessionId  string
	startAt        time.Time
	chatService    public.ChatService
	sentWsMessages uint32
	ws             Ws
	sync.RWMutex
}

func newSession(userId, username, chatUserId string, chatSessionId string, chatService public.ChatService) Session {

	s := &sessionImpl{
		id:            kit.NewId(),
		userId:        userId,
		username:      username,
		chatUserId:    chatUserId,
		chatSessionId: chatSessionId,
		chatService:   chatService,
		startAt:       time.Now().UTC(),
	}

	return s
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

func (s *sessionImpl) wsListen(ws Ws) {

	closedEventChan := ws.closedEvent()
	receivedEventChan := ws.receivedMessageEvent()
	for {
		select {
		case <-closedEventChan:
			s.Lock()
			s.ws = nil
			s.Unlock()
			return
		case msg := <-receivedEventChan:
			//TODO: here we have to forward message to Chat service if we need
			log.TrcF("message forward to chat: %v", msg)
		}
	}
}

func (s *sessionImpl) sendWsMessage(msg []byte) error {

	log.TrcF("[session][ws] message sending to user=%s msg=%s", s.userId, string(msg))

	// check if WS connection is open
	if s.isWs() {
		s.ws.send(msg)
		// increment sent message counter
		atomic.AddUint32(&s.sentWsMessages, 1)
	} else {
		log.InfF("[session][ws] cannot send WS message. no client connection")
	}
	return nil
}

func (s *sessionImpl) setWs(ws Ws) Session {

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

func (s *sessionImpl) setJWT(jwt *auth.JWT) Session {
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

	if s.chatSessionId != "" {
		go func() {
			err := s.chatService.Logout(ctx, s.chatUserId)
			if err != nil {
				log.Err(err, true)
			}
		}()
	}

	if s.ws != nil {
		s.ws.close()
	}

}
