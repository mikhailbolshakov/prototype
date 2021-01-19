package infrastructure

import (
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
)

type Container struct {
	Bpm   bpm.Engine
}

func New() *Container {
	c := &Container{
		Bpm: zeebe.NewEngine(&zeebe.Params{
			Port: "26500",
			Host: "localhost",
		}),
	}
	return c
}

func (c *Container) Init() error {
	return c.Bpm.Open()
}

func Close() {}
