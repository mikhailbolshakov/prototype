package session

import (
	"encoding/json"
	"fmt"
	"github.com/adacta-ru/mattermost-server/v6/model"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/chat/mattermost"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"sync"
)

type Session interface {
	setWs(ws Ws) Session
	setJWT(jwt *auth.JWT) Session
	isWs() bool
	getId() string
	getUserId() string
	getUsername() string
	getAccessToken() string
	close()
}

type sessionImpl struct {
	Id           string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	UserId       string
	Username     string
	mmClient     *mattermost.Client
	ws           Ws
	sync.RWMutex
	closeMMws chan struct{}
}

func newSession(userId, username string, mmClient *mattermost.Client) Session {
	s := &sessionImpl{
		Id:        kit.NewId(),
		UserId:    userId,
		Username:  username,
		mmClient:  mmClient,
		closeMMws: make(chan struct{}),
	}

	// start listening MM ws
	if s.mmClient != nil {
		s.listenMMSocket()
	}

	return s
}

func (s *sessionImpl) forwardEventFromMM(event *model.WebSocketEvent) {

	log.TrcF("[mm-ws] event=%s data=%v", event.EventType(), event.ToJson())

	s.RLock()
	defer s.RUnlock()

	if s.ws != nil {
		msg, _ := json.Marshal(event)
		s.ws.send(msg)
	}

}

func (s *sessionImpl) forwardResponseFromMM(response *model.WebSocketResponse) {
	log.TrcF("[mm-ws] event=%s data=%v", response.EventType(), response.ToJson())

	s.RLock()
	defer s.RUnlock()

	if s.ws != nil {
		msg, _ := json.Marshal(response)
		s.ws.send(msg)
	}

}

func (s *sessionImpl) listenMMSocket() {

	go s.mmClient.WsApi.Listen()
	go func() {
		for {
			select {
			case event := <-s.mmClient.WsApi.EventChannel:
				s.forwardEventFromMM(event)
			case response := <-s.mmClient.WsApi.ResponseChannel:
				s.forwardResponseFromMM(response)
			case _ = <-s.mmClient.WsApi.PingTimeoutChannel:
				log.Trc("[mm-ws] ping")
			case <-s.closeMMws:
				log.Trc("[mm-ws] close", s.mmClient.User.Email)
				s.mmClient.WsApi.Close()
				s.mmClient.RestApi.ClearOAuthToken()
				return
			}
		}
	}()
}

func (s *sessionImpl) getId() string {
	s.RLock()
	defer s.RUnlock()
	return s.Id
}

func (s *sessionImpl) getUserId() string {
	s.RLock()
	defer s.RUnlock()
	return s.UserId
}

func (s *sessionImpl) getUsername() string {
	s.RLock()
	defer s.RUnlock()
	return s.Username
}

func (s *sessionImpl) getAccessToken() string {
	s.RLock()
	defer s.RUnlock()
	return s.AccessToken
}

func (s *sessionImpl) forwardToMM(msg []byte) {

	if s.mmClient == nil {
		log.Dbg("[ws] received message to MM, but mm client is nil")
		return
	}

	var m map[string]interface{}
	if err := json.Unmarshal(msg, &m); err != nil {
		log.Err(err, true)
		return
	}

	action, ok := m["action"]
	if !ok {
		log.Err(fmt.Errorf("[ws] message to MM must have 'action'"), false)
		return
	}

	if data, ok := m["data"]; !ok {
		log.Err(fmt.Errorf("[ws] message to MM must have 'data'"), false)
		return
	} else {
		s.mmClient.WsApi.SendMessage(action.(string), data.(map[string]interface{}))
		log.TrcF("[ws] resent to MM: %s", string(msg))
	}

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
			//TODO: we have to distinguish between MM message and ours
			s.forwardToMM(msg)
		}
	}
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
	s.Lock()
	defer s.Unlock()
	s.AccessToken = jwt.AccessToken
	s.RefreshToken = jwt.RefreshToken
	s.ExpiresIn = jwt.ExpiresIn
	return s
}

func (s *sessionImpl) isWs() bool {
	s.RLock()
	defer s.RUnlock()
	return s.ws != nil
}

func (s *sessionImpl) close() {

	s.Lock()
	defer s.Unlock()

	if s.mmClient != nil {
		s.closeMMws <- struct{}{}
	}

	if s.ws != nil {
		s.ws.close()
	}

}
