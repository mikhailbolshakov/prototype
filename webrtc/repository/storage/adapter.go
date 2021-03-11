package storage

import (
	kitCache "gitlab.medzdrav.ru/prototype/kit/cache"
	"gitlab.medzdrav.ru/prototype/kit/kv"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/logger"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() domain.RoomStorage
	GetRoomCoordinator() domain.RoomCoordinator
	Close()
}

type container struct {
	Db         *kitStorage.Storage
	ReadOnlyDB *kitStorage.Storage
	Cache      *kitCache.Redis
	Etcd       *kv.Etcd
}

type adapterImpl struct {
	container     *container
	storageImpl   *storageImpl
	roomCoordImpl *etcdRoomCoordImpl
}

func NewAdapter() Adapter {
	a := &adapterImpl{}
	a.container = &container{}
	a.storageImpl = newStorage(a.container)
	a.roomCoordImpl = &etcdRoomCoordImpl{ c: a.container }
	return a
}

func (a *adapterImpl) Init(cfg *config.Config) error {

	servCfg := cfg.Services["webrtc"]

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

	a.container.Etcd, err = kv.Open(&kv.Options{Hosts: cfg.Etcd.Hosts}, logger.LF())
	if err != nil {
		return err
	}

	return nil
}

func (a *adapterImpl) GetService() domain.RoomStorage {
	return a.storageImpl
}

func (a *adapterImpl) GetRoomCoordinator() domain.RoomCoordinator {
	return a.roomCoordImpl
}

func (a *adapterImpl) Close() {
	a.container.Db.Close()
	a.container.ReadOnlyDB.Close()
	a.container.Cache.Close()
	_ = a.container.Etcd.Close()
}
