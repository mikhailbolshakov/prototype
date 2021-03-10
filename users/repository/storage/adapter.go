package storage

import (
	kitCache "gitlab.medzdrav.ru/prototype/kit/cache"
	"gitlab.medzdrav.ru/prototype/kit/search"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/proto/config"
	domain "gitlab.medzdrav.ru/prototype/users/domain"
	"gitlab.medzdrav.ru/prototype/users/logger"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() domain.UserStorage
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

func (a *adapterImpl) Init(cfg *config.Config) error {

	servCfg := cfg.Services["users"]

	var err error

	// storage R/W
	a.container.Db, err = kitStorage.Open(&kitStorage.Params{
		UserName: servCfg.Database.User,
		Password: servCfg.Database.Password,
		DBName:   servCfg.Database.Dbname,
		Port:     servCfg.Database.Port,
		Host:     servCfg.Database.HostRw,
	}, logger.LF())
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
	}, logger.LF())
	if err != nil {
		return err
	}

	// Redis
	a.container.Cache, err = kitCache.Open(&kitCache.Params{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		Ttl:      uint(cfg.Redis.Ttl),
	}, logger.LF())
	if err != nil {
		return err
	}

	// Index search
	a.container.Search, err = search.NewEs(cfg.Es.Url, cfg.Es.Trace, logger.LF())
	if err != nil {
		return err
	}
	err = a.storageImpl.ensureIndex()
	if err != nil {
		return err
	}

	return nil
}

func (a *adapterImpl) GetService() domain.UserStorage {
	return a.storageImpl
}

func (a *adapterImpl) Close() {
	a.container.Db.Close()
	a.container.ReadOnlyDB.Close()
	a.container.Cache.Close()
	a.container.Search.Close()
}
