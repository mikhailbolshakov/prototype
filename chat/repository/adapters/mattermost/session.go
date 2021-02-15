package mattermost

import (
	"context"
	"fmt"
	"github.com/adacta-ru/mattermost-server/v6/model"
	"gitlab.medzdrav.ru/prototype/kit"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"sync"
)

type ChatSession interface {
	GetId() string
	GetUserId() string
	GetChatUserId() string
	GetUsername() string
	Open(ctx context.Context)
	Close(ctx context.Context)
	Client() *Client
}

type sessionImpl struct {
	Id         string
	UserId     string
	ChatUserId string
	Username   string
	mmClient   *Client
	sync.RWMutex
	closeMMws chan struct{}
	queue     queue.Queue
}

func newSession(ctx context.Context, userId, username, chatUserId string, cfg *kitConfig.Config, queue queue.Queue) (ChatSession, error) {

	s := &sessionImpl{
		Id:         kit.NewId(),
		UserId:     userId,
		Username:   username,
		ChatUserId: chatUserId,
		closeMMws:  make(chan struct{}),
		queue:      queue,
	}

	mmClient, err := Login(&Params{
		Url:     cfg.Mattermost.Url,
		WsUrl:   cfg.Mattermost.Ws,
		Account: username,
		// currently we don't care about chat users' passwords
		Pwd:    cfg.Mattermost.DefaultPassword,
		OpenWS: true,
	})
	if err != nil {
		return nil, err
	}
	s.mmClient = mmClient

	return s, nil
}

func newAdminSession(ctx context.Context, cfg *kitConfig.Config) (ChatSession, error) {

	s := &sessionImpl{
		Id: kit.NewId(),
	}

	mmClient, err := Login(&Params{
		Url:     cfg.Mattermost.Url,
		WsUrl:   cfg.Mattermost.Ws,
		Account: cfg.Mattermost.AdminUsername,
		Pwd:     cfg.Mattermost.AdminPassword,
		OpenWS:  false,
	})
	if err != nil {
		return nil, err
	}
	s.mmClient = mmClient

	return s, nil
}

func newBotSession(ctx context.Context, cfg *kitConfig.Config) (ChatSession, error) {

	s := &sessionImpl{
		Id: kit.NewId(),
	}

	mmClient, err := Login(&Params{
		Url:         cfg.Mattermost.Url,
		WsUrl:       cfg.Mattermost.Ws,
		Account:     cfg.Mattermost.BotUsername,
		AccessToken: cfg.Mattermost.BotAccessToken,
		OpenWS:      false,
	})
	if err != nil {
		return nil, err
	}
	s.mmClient = mmClient

	return s, nil
}

func (s *sessionImpl) skipEvent(eventType string) bool {
	_, ok := skippedWsEventTypes[eventType]
	return ok
}

func (s *sessionImpl) forwardEvent(ctx context.Context, event *model.WebSocketEvent) {

	log.TrcF("[mm-ws] event=%s data=%v", event.EventType(), event.ToJson())

	if s.skipEvent(event.EventType()) {
		log.TrcF("[mm-ws] event=%s skipped")
		return
	}

	s.RLock()
	userId := s.UserId
	chatUserId := s.ChatUserId
	s.RUnlock()

	c, ok := kitContext.Request(ctx)
	if !ok {
		log.Err(fmt.Errorf("invalid context"), true)
		return
	}

	msg := &queue.Message{
		Ctx: c,
		Payload: &ChatIncomingWsEventPayload{
			UserId:     userId,
			ChatUserId: chatUserId,
			Event:      event.EventType(),
			Data:       event.GetData(),
		},
	}

	if err := s.queue.Publish(ctx, queue.QUEUE_TYPE_AT_MOST_ONCE, INCOMING_WS_QUEUE_TOPIC, msg); err != nil {
		log.Err(err, true)
		return
	}

}

func (s *sessionImpl) listen(ctx context.Context) {

	if !s.mmClient.Params.OpenWS {
		return
	}

	go s.mmClient.WsApi.Listen()
	go func() {
		for {
			select {
			case event := <-s.mmClient.WsApi.EventChannel:
				s.forwardEvent(ctx, event)
			case response := <-s.mmClient.WsApi.ResponseChannel:
				log.TrcF("[mm-ws] response event=%s data=%s", response.EventType(), response.ToJson())
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

func (s *sessionImpl) GetId() string {
	s.RLock()
	defer s.RUnlock()
	return s.Id
}

func (s *sessionImpl) Client() *Client {
	s.RLock()
	defer s.RUnlock()
	return s.mmClient
}

func (s *sessionImpl) GetUserId() string {
	s.RLock()
	defer s.RUnlock()
	return s.UserId
}

func (s *sessionImpl) GetChatUserId() string {
	s.RLock()
	defer s.RUnlock()
	return s.ChatUserId
}

func (s *sessionImpl) GetUsername() string {
	s.RLock()
	defer s.RUnlock()
	return s.Username
}

func (s *sessionImpl) Open(ctx context.Context) {
	s.listen(ctx)
}

func (s *sessionImpl) Close(ctx context.Context) {

	s.Lock()
	defer s.Unlock()

	if s.mmClient.Params.OpenWS {
		s.closeMMws <- struct{}{}
	} else {
		s.mmClient.RestApi.ClearOAuthToken()
	}

}
