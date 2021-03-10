package mattermost

import (
	"context"
	"github.com/adacta-ru/mattermost-server/v6/model"
	"gitlab.medzdrav.ru/prototype/chat/logger"
	"gitlab.medzdrav.ru/prototype/kit"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"sync"
)

// TODO: move to config
// skippedWsEventTypes defines mattermost events to skip
// the whole list of available events find at https://api.mattermost.com/#tag/WebSocket
var skippedWsEventTypes = map[string]struct{}{
	"hello":                   {},
	"config_changed":          {},
	"delete_team":             {},
	"leave_team":              {},
	"license_changed":         {},
	"plugin_disabled":         {},
	"plugin_enabled":          {},
	"plugin_statuses_changed": {},
	"preference_changed":      {},
	"preferences_changed":     {},
	"preferences_deleted":     {},
	"response":                {},
	"update_team":             {},
	"user_added":              {},
	"user_removed":            {},
	"user_role_updated":       {},
	"user_updated":            {},
}

type mattermostWsEvent struct {
	ChatUserId string                 `json:"chatUserId"`
	Event      string                 `json:"event"`
	Data       map[string]interface{} `json:"data"`
}

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

func newSession(ctx context.Context, userId, username, chatUserId string, cfg *config.Config, queue queue.Queue) (ChatSession, error) {

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

func newAdminSession(ctx context.Context, cfg *config.Config) (ChatSession, error) {

	s := &sessionImpl{
		Id: kit.NewId(),
	}

	mmClient, err := Login(&Params{
		Url:         cfg.Mattermost.Url,
		WsUrl:       cfg.Mattermost.Ws,
		Account:     cfg.Mattermost.AdminUsername,
		AccessToken: cfg.Mattermost.AdminAccessToken,
		OpenWS:      false,
	})
	if err != nil {
		return nil, err
	}
	s.mmClient = mmClient

	return s, nil
}

func newBotSession(ctx context.Context, cfg *config.Config) (ChatSession, error) {

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

func (s *sessionImpl) l() log.CLogger {
	return logger.L().Cmp("mm-session")
}

func (s *sessionImpl) forwardEvent(event *model.WebSocketEvent) {

	l := s.l().Mth("fwd-ev").F(log.FF{"type": event.EventType()})

	l.Dbg().Trc(event.ToJson())

	if s.skipEvent(event.EventType()) {
		l.Trc("skipped")
		return
	}

	s.RLock()
	userId := s.UserId
	chatUserId := s.ChatUserId
	username := s.Username
	s.RUnlock()

	ctx := kitContext.NewRequestCtx().
		Client("mm-ws").
		WithNewRequestId().
		WithUser(userId, username).
		WithChatUserId(chatUserId)

	msg := &queue.Message{
		Ctx: ctx,
		Payload: &proto.OutgoingWsEventQueueMessagePayload{
			UserId: userId,
			WsEvent: &proto.WsMessage{
				MessageType: proto.WS_MESSAGE_TYPE_CHAT,
				Data: &mattermostWsEvent{
					ChatUserId: chatUserId,
					Event:      event.EventType(),
					Data:       event.GetData(),
				},
			},
		},
	}

	if err := s.queue.Publish(ctx.ToContext(context.Background()), queue.QUEUE_TYPE_AT_MOST_ONCE, proto.QUEUE_TOPIC_OUTGOING_WS_EVENT, msg); err != nil {
		l.E(err).Err("publish failed")
		return
	}

	l.Dbg("ok")

}

func (s *sessionImpl) listen(ctx context.Context) {

	if !s.mmClient.Params.OpenWS {
		return
	}

	go s.mmClient.WsApi.Listen()
	go func() {

		l := s.l().Mth("listen")

		for {
			select {
			case event := <-s.mmClient.WsApi.EventChannel:
				s.forwardEvent(event)
			case response := <-s.mmClient.WsApi.ResponseChannel:
				l.F(log.FF{"event": response.EventType()}).Dbg().Trc(response.ToJson())
			case _ = <-s.mmClient.WsApi.PingTimeoutChannel:
				l.Trc("ping")
			case <-s.closeMMws:
				s.mmClient.WsApi.Close()
				s.mmClient.RestApi.ClearOAuthToken()
				l.TrcF("%s closed\n", s.mmClient.User.Username)
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
