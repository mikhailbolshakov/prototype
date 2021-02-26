package storage

import (
	kitCache "gitlab.medzdrav.ru/prototype/kit/cache"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
	GetService() domain.WebrtcStorage
	Close()
}

type container struct {
	Db         *kitStorage.Storage
	ReadOnlyDB *kitStorage.Storage
	Cache      *kitCache.Redis
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

	servCfg := cfg.Services["webrtc"]

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

	return nil
}

func (a *adapterImpl) GetService() domain.WebrtcStorage {
	return a.storageImpl
}

func (a *adapterImpl) Close() {
	a.container.Db.Close()
	a.container.ReadOnlyDB.Close()
	a.container.Cache.Close()
}
