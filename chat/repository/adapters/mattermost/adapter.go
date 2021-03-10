package mattermost

import (
	"context"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/chat/logger"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"time"
)

type Adapter interface {
	Init(ctx context.Context, cfg *config.Config) error
	GetService() domain.MattermostService
	Close(ctx context.Context)
}

type adapterImpl struct {
	mmServiceImpl *serviceImpl
	sessionHub    ChatSessionHub
	quit          chan struct{}
}

func NewAdapter(hub ChatSessionHub) Adapter {
	a := &adapterImpl{
		sessionHub:    hub,
		mmServiceImpl: newImpl(hub),
		quit: make(chan struct{}),
	}
	return a
}

// MattermostKeepAlive pings Mattermost periodically (1s by default) and if it's not available set "notReady" flag for the service, so that it won't able to handle requests
// As Mattermost is up, it initializes connections
func (a *adapterImpl) MattermostKeepAlive(ctx context.Context) {

	url := a.mmServiceImpl.cfg.Mattermost.Url
	at := a.mmServiceImpl.cfg.Mattermost.AdminAccessToken

	go func() {

		l := logger.L().Cmp("mm").Mth("keepalive")

		wasReady := false
		for {
			select {
			case <-time.NewTicker(time.Second).C:
				if ping(url, at) {
					if !wasReady {
						if err := a.sessionHub.Init(ctx); err != nil {
							l.E(err).Err("hub init failed")
							break
						}
						a.mmServiceImpl.setReady(true)
						wasReady = true
						l.Trc("ok")
					}
				} else {
					wasReady = false
					a.mmServiceImpl.setReady(false)
					l.Trc("readiness prob failed")
				}
			case <-a.quit:
				return
			}
		}

	}()
}

func (a *adapterImpl) Init(ctx context.Context, cfg *config.Config) error {
	a.mmServiceImpl.setConfig(cfg)
	a.MattermostKeepAlive(ctx)
	return nil
}

func (a *adapterImpl) GetService() domain.MattermostService {
	return a.mmServiceImpl
}

func (a *adapterImpl) Close(ctx context.Context) {
	close(a.quit)
	a.sessionHub.Close(ctx)
}
