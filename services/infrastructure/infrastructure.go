package infrastructure

import (
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	kitCache "gitlab.medzdrav.ru/prototype/kit/cache"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
)

type Container struct {
	Db         *kitStorage.Storage
	ReadOnlyDB *kitStorage.Storage
	Cache      *kitCache.Redis
	Bpm        bpm.Engine
}

func New() *Container {
	c := &Container{}
	return c
}

func (c *Container) Init(cfg *kitConfig.Config) error {

	servCfg := cfg.Services["services"]

	var err error

	// storage R/W
	c.Db, err = kitStorage.Open(&kitStorage.Params{
		UserName: servCfg.Database.User,
		Password: servCfg.Database.Password,
		DBName:   servCfg.Database.Dbname,
		Port:     servCfg.Database.Port,
		Host:     servCfg.Database.HostRw,
	})
	if err != nil {
		return err
	}

	// storage Readonly
	c.ReadOnlyDB, err = kitStorage.Open(&kitStorage.Params{
		UserName: servCfg.Database.User,
		Password: servCfg.Database.Password,
		DBName:   servCfg.Database.Dbname,
		Port:     servCfg.Database.Port,
		Host:     servCfg.Database.HostRo,
	})
	if err != nil {
		return err
	}

	// Redis
	c.Cache, err = kitCache.Open(&kitCache.Params{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		Ttl:      uint(cfg.Redis.Ttl),
	})
	if err != nil {
		return err
	}

	c.Bpm = zeebe.NewEngine(&zeebe.Params{
		Port: cfg.Zeebe.Port,
		Host: cfg.Zeebe.Host,
	})

	if err := c.Bpm.Open(); err != nil {
		return err
	}

	return nil
}

func (c *Container) Close() {
	c.Db.Close()
	c.ReadOnlyDB.Close()
	_ = c.Cache.Instance.Close()
	_ = c.Bpm.Close()
}
