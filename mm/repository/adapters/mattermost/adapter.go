package mattermost

import "gitlab.medzdrav.ru/prototype/kit/queue"

type Adapter interface {
	Init() error
	GetService() Service
	ListenAsync() error
}

type adapterImpl struct {
	mmServiceImpl *serviceImpl
	initialized bool
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
	a.mmServiceImpl.client, err = login(&Params{
		Url:     "http://localhost:8065",
		WsUrl:   "ws://localhost:8065",
		Account: "admin",
		Pwd:     "admin",
	})
	if err != nil {
		panic(err)
	}

	err = a.mmServiceImpl.queue.Open("prototype-mm")
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