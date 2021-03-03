package metrcics

import (
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/sessions/domain"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
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

func (a *adapterImpl) Init(c *kitConfig.Config) error {
	return nil
}

func (a *adapterImpl) GetService() domain.Metrics {
	return a.serviceImpl
}

func (a *adapterImpl) Close()  {
}
