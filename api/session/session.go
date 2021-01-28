package session

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/chat/mattermost"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"log"
	"sync"
)

type Session interface {
	setWs(ws Ws) Session
	setJWT(jwt *auth.JWT) Session
	isWs() bool
	getId() string
	getAccessToken() string
	close()
}

type sessionImpl struct {
	Id           string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	UserId       string
	mmClient     *mattermost.Client
	ws           Ws
	sync.RWMutex
}

func newSession(userId string, mmClient *mattermost.Client) Session {
	s := &sessionImpl{
		Id:       kit.NewId(),
		UserId:   userId,
		mmClient: mmClient,
	}

	if mmClient != nil {
		s.listenMMSocket()
	}

	return s
}

func (s *sessionImpl) listenMMSocket() {

	go s.mmClient.WsApi.Listen()
	go func() {
		for {
			select {
			case event := <-s.mmClient.WsApi.EventChannel:
				s, _ := json.MarshalIndent(event, "", "\t")
				log.Printf("[WS event]. %s", s)
			case response := <-s.mmClient.WsApi.ResponseChannel:
				s, _ := json.MarshalIndent(response, "", "\t")
				log.Printf("[WS response]. %s", s)
			case <-s.mmClient.Quit:
				log.Printf("[WS closing]. Closing request for user %s", s.mmClient.User.Email)
				s.mmClient.WsApi.Close()
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

func (s *sessionImpl) getAccessToken() string {
	s.RLock()
	defer s.RUnlock()
	return s.AccessToken
}

func (s *sessionImpl) setWs(ws Ws) Session {

	s.Lock()
	defer s.Unlock()

	if s.ws != nil {
		s.ws.close()
	}

	s.ws = ws
	ws.open()

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
		s.mmClient.WsApi.Close()
		s.mmClient.RestApi.ClearOAuthToken()
		s.mmClient.Quit <- struct{}{}
	}

	if s.ws != nil {
		s.ws.close()
	}

}