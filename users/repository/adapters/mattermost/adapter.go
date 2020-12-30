package mattermost

import (
	"gitlab.medzdrav.ru/prototype/mm"
)

type Adapter interface {
	Init() error
	GetService() Service
	ListenAsync() error
}

type adapterImpl struct {
	mmServiceImpl *serviceImpl
	initialized bool
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		mmServiceImpl: newImpl(),
		initialized:   false,
	}
	return a
}

func (a *adapterImpl) Init() error {

	var err error
	a.mmServiceImpl.client, err = mm.Login(&mm.Params{
		Url:     "http://localhost:8065",
		WsUrl:   "ws://localhost:8065",
		Account: "admin",
		Pwd:     "admin",
	})
	if err != nil {
		panic(err)
	}

	a.mmServiceImpl.queue = mm.NewQueueHandler()
	err = a.mmServiceImpl.queue.Open("prototype-tasks")
	if err != nil {
		return err
	}
	a.initialized = true

	return nil
}

func (a *adapterImpl) GetService() Service {
	return a.mmServiceImpl
}

func (a *adapterImpl) ListenAsync() error {
	return a.mmServiceImpl.listenNewPosts()
}