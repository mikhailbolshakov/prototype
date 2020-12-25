package repository

import (
	kitCache "gitlab.medzdrav.ru/prototype/kit/cache"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/mm"
)

// TODO: I don't like it
var storage *kitStorage.Storage
var cache *kitCache.Redis
var mmClient *mm.Client

func InitInfrastructure() error {

	var err error

	// storage
	storage, err = kitStorage.Open(&kitStorage.Params{
		UserName: "users",
		Password: "users",
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

	mmClient, err = mm.Login(&mm.Params{
		Url:     "http://localhost:8065",
		WsUrl:   "ws://localhost:8065",
		Account: "admin",
		Pwd:     "admin",
	})
	if err != nil {
		return err
	}

	return nil
}

func CloseInfrastructure() {

}