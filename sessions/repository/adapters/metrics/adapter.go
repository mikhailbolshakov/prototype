package metrcics

import (
	"gitlab.medzdrav.ru/prototype/proto/config"
	"gitlab.medzdrav.ru/prototype/sessions/domain"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() domain.Metrics
	Close()
}

type adapterImpl struct {
	serviceImpl *serviceImpl
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		serviceImpl: newImpl(),
	}
	return a
}

func (a *adapterImpl) Init(c *config.Config) error {
	return nil
}

func (a *adapterImpl) GetService() domain.Metrics {
	return a.serviceImpl
}

func (a *adapterImpl) Close()  {
}
