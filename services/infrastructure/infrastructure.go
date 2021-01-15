package infrastructure

import (
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	kitCache "gitlab.medzdrav.ru/prototype/kit/cache"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
)

type Container struct {
	Db    *kitStorage.Storage
	Cache *kitCache.Redis
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

	var err error

	// storage
	c.Db, err = kitStorage.Open(&kitStorage.Params{
		UserName: "services",
		Password: "services",
		DBName:   "mattermost",
		Port:     "5432",
		Host:     "localhost",
	})
	if err != nil {
		return err
	}

	// Redis
	c.Cache, err = kitCache.Open(&kitCache.Params{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		Ttl:      7200,
	})
	if err != nil {
		return err
	}

	if err := c.Bpm.Open(); err != nil {
		return err
	}

	return nil
}

func Close() {}
