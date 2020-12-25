package repository

import (
	kitCache "gitlab.medzdrav.ru/prototype/kit/cache"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
)

// TODO: I don't like it
var storage *kitStorage.Storage
var cache *kitCache.Redis

func InitInfrastructure() error {

	var err error

	// storage
	storage, err = kitStorage.Open(&kitStorage.Params{
		UserName: "tasks",
		Password: "tasks",
		DBName:   "mattermost",
		Port:     "5432",
		Host:     "localhost",
	})
	if err != nil {
		return err
	}

	// Redis
	cache, err = kitCache.Open(&kitCache.Params{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		Ttl:      7200,
	})
	if err != nil {
		return err
	}

	return nil
}

func CloseInfrastructure() {

}