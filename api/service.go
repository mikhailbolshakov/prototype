package api

import (
	"gitlab.medzdrav.ru/prototype/api/tasks"
	"gitlab.medzdrav.ru/prototype/api/users"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"log"
)

type Service struct {
	*kitHttp.Server
}

func NewService() *Service {

	return &Service{
		// TODO: take from .env
		Server: kitHttp.NewHttpServer("localhost", "8000", &users.Router{}, &tasks.Router{}),
	}
}

func (u *Service) Start() error {

	go func() {
		log.Fatal(u.Open())
	}()

	return nil
}
