package api

import (
	"gitlab.medzdrav.ru/prototype/api/tasks"
	"gitlab.medzdrav.ru/prototype/api/users"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"log"
)

type serviceImpl struct {
	*kitHttp.Server
}

func New() service.Service {
	return &serviceImpl{}
}

func (u *serviceImpl) Init() error {
	u.Server = kitHttp.NewHttpServer("localhost", "8000", &users.Router{}, &tasks.Router{})
	return nil
}

func (u *serviceImpl) Listen() error {
	return nil
}

func (u *serviceImpl) ListenAsync() error {

	go func() {
		log.Fatal(u.Open())
	}()


	return nil
}
