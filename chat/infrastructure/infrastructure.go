package infrastructure

import kitConfig "gitlab.medzdrav.ru/prototype/kit/config"

type Container struct {
}

func New() *Container {
	c := &Container{
	}
	return c
}

func (c *Container) Init(cfg *kitConfig.Config) error {
	return nil
}

func (c *Container) Close() {
}
