package mattermost

import (
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/kit/chat/mattermost"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
)

type Adapter interface {
	Init(cfg *kitConfig.Config) error
	GetService() domain.MattermostService
	Close()
}

type adapterImpl struct {
	mmServiceImpl *serviceImpl
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		mmServiceImpl: newImpl(),
	}
	return a
}

func (a *adapterImpl) Init(cfg *kitConfig.Config) error {

	a.mmServiceImpl.setConfig(cfg)

	var err error
	a.mmServiceImpl.adminClient, err = mattermost.Login(&mattermost.Params{
		Url:     cfg.Mattermost.Url,
		WsUrl:   cfg.Mattermost.Ws,
		Account: cfg.Mattermost.AdminUsername,
		Pwd:     cfg.Mattermost.AdminPassword,
	})
	if err != nil {
		return err
	}

	a.mmServiceImpl.botClient, err = mattermost.Login(&mattermost.Params{
		Account:     cfg.Mattermost.BotUsername,
		Url:         cfg.Mattermost.Url,
		AccessToken: cfg.Mattermost.BotAccessToken,
	})
	return err

}

func (a *adapterImpl) GetService() domain.MattermostService {
	return a.mmServiceImpl
}

func (a *adapterImpl) Close() {
	if a.mmServiceImpl.adminClient.WsApi != nil {
		a.mmServiceImpl.adminClient.WsApi.Close()
	}
	a.mmServiceImpl.adminClient.RestApi.ClearOAuthToken()

	if a.mmServiceImpl.botClient.WsApi != nil {
		a.mmServiceImpl.botClient.WsApi.Close()
	}
	a.mmServiceImpl.botClient.RestApi.ClearOAuthToken()
}
