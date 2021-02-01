package mattermost

import (
	"gitlab.medzdrav.ru/prototype/kit/chat/mattermost"
	"gitlab.medzdrav.ru/prototype/kit/queue"
)

type Adapter interface {
	Init() error
	GetService() Service
	Close()
}

type adapterImpl struct {
	mmServiceImpl *serviceImpl
	initialized   bool
}

func NewAdapter(queue queue.Queue) Adapter {
	a := &adapterImpl{
		mmServiceImpl: newImpl(queue),
		initialized:   false,
	}
	return a
}

func (a *adapterImpl) Init() error {

	var err error
	a.mmServiceImpl.adminClient, err = mattermost.Login(&mattermost.Params{
		Url:     MM_REST_URL,
		WsUrl:   MM_WS_URL,
		Account: ADMIN_USERNAME,
		Pwd:     ADMIN_PASSWORD,
	})
	if err != nil {
		return err
	}

	a.mmServiceImpl.botClient, err = mattermost.Login(&mattermost.Params{
		Account:     RGS_BOT_USERNAME,
		Url:         MM_REST_URL,
		AccessToken: RGS_BOT_ACCESS_TOKEN,
	})
	return err

}

func (a *adapterImpl) GetService() Service {
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
