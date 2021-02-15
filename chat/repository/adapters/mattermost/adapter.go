package mattermost

import (
	"context"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
)

type Adapter interface {
	Init(ctx context.Context, cfg *kitConfig.Config) error
	GetService() domain.MattermostService
	Close(ctx context.Context)
}

type adapterImpl struct {
	mmServiceImpl *serviceImpl
	sessionHub    ChatSessionHub
}

func NewAdapter(hub ChatSessionHub) Adapter {
	a := &adapterImpl{
		sessionHub: hub,
		mmServiceImpl: newImpl(hub),
	}
	return a
}

func (a *adapterImpl) Init(ctx context.Context, cfg *kitConfig.Config) error {

	if err := a.sessionHub.Init(ctx); err != nil {
		return err
	}

	a.mmServiceImpl.setConfig(cfg)

	return nil

}

func (a *adapterImpl) GetService() domain.MattermostService {
	return a.mmServiceImpl
}

func (a *adapterImpl) Close(ctx context.Context) {
	a.sessionHub.Close(ctx)
}
