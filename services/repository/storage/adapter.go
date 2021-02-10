package storage

import (
	kitCache "gitlab.medzdrav.ru/prototype/kit/cache"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/search"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/services/domain"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
	GetService() domain.Storage
	Close()
}

type container struct {
	Db         *kitStorage.Storage
	ReadOnlyDB *kitStorage.Storage
	Cache      *kitCache.Redis
	Search     search.Search
}

type adapterImpl struct {
	container   *container
	storageImpl *storageImpl
}

func NewAdapter() Adapter {
	a := &adapterImpl{}
	a.container = &container{}
	a.storageImpl = newStorage(a.container)
	return a
}

func (a *adapterImpl) Init(cfg *kitConfig.Config) error {

	servCfg := cfg.Services["services"]

	var err error

	// storage R/W
	a.container.Db, err = kitStorage.Open(&kitStorage.Params{
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
	a.container.ReadOnlyDB, err = kitStorage.Open(&kitStorage.Params{
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
	a.container.Cache, err = kitCache.Open(&kitCache.Params{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		Ttl:      uint(cfg.Redis.Ttl),
	})
	if err != nil {
		return err
	}

	// Index search
	a.container.Search, err = search.NewEs(cfg.Es.Url, cfg.Es.Trace)
	if err != nil {
		return err
	}
	err = a.storageImpl.ensureIndex()
	if err != nil {
		return err
	}

	return nil
}

func (a *adapterImpl) GetService() domain.Storage {
	return a.storageImpl
}

func (a *adapterImpl) Close() {
	a.container.Db.Close()
	a.container.ReadOnlyDB.Close()
	a.container.Cache.Close()
	a.container.Search.Close()
}
