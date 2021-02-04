package infrastructure

import (
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
)

type Container struct {
	Bpm bpm.Engine
}

func New() *Container {
	c := &Container{}
	return c
}

func (c *Container) Init(cfg *kitConfig.Config) error {
	c.Bpm = zeebe.NewEngine(&zeebe.Params{
		Port: cfg.Zeebe.Port,
		Host: cfg.Zeebe.Host,
	})
	return c.Bpm.Open()
}

func (c *Container) Close() {
	_ = c.Bpm.Close()
}
