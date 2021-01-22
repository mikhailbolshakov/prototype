package api

import (
	"context"
	"github.com/Nerzal/gocloak/v7"
	"gitlab.medzdrav.ru/prototype/api/services"
	"gitlab.medzdrav.ru/prototype/api/tasks"
	"gitlab.medzdrav.ru/prototype/api/users"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"log"
)

type serviceImpl struct {
	keycloak gocloak.GoCloak
	mdw      auth.AuthMiddleware
	*kitHttp.Server
}

func New() service.Service {
	return &serviceImpl{}
}

func (u *serviceImpl) Init() error {

	ctx := context.Background()

	authClient := &auth.AuthClient{
		ClientID:     "app",
		ClientSecret: "d6dbae97-8570-4758-a081-9077b7899a7d",
		Realm:        "prototype",
	}

	u.keycloak = gocloak.NewClient("http://localhost:8086")
	u.mdw = auth.NewKeyCloakMdw(ctx, u.keycloak, authClient, "", "")

	authHandler := auth.NewAuthenticationHandler(ctx, u.keycloak, authClient)
	u.Server = kitHttp.NewHttpServer("localhost", "8000", users.New(authHandler), tasks.New(), services.New())

	// set auth middleware
	u.Server.SetAuthMiddleware(u.mdw.CheckToken)

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
