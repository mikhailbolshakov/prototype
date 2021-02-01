package api

import (
	"context"
	"github.com/Nerzal/gocloak/v7"
	"gitlab.medzdrav.ru/prototype/api/public/bp"
	servRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/services"
	tasksRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/tasks"
	usersRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/api/public/services"
	"gitlab.medzdrav.ru/prototype/api/session"
	"gitlab.medzdrav.ru/prototype/api/public/tasks"
	"gitlab.medzdrav.ru/prototype/api/public/users"
	bpRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/bp"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"log"
)

type serviceImpl struct {
	*kitHttp.Server
	keycloak    gocloak.GoCloak
	mdw         auth.Middleware
	hub         session.Hub
	userAdapter usersRep.Adapter
	taskAdapter tasksRep.Adapter
	servAdapter servRep.Adapter
	bpAdapter   bpRep.Adapter
}

func New() service.Service {
	s := &serviceImpl{}

	ctx := context.Background()

	authClient := &auth.ClientSecurityInfo{
		ID:     "app",
		Secret: "d6dbae97-8570-4758-a081-9077b7899a7d",
		Realm:  "prototype",
	}

	s.keycloak = gocloak.NewClient("http://localhost:8086")
	s.mdw = auth.NewMdw(ctx, s.keycloak, authClient, "", "")

	authService := auth.New(ctx, s.keycloak, authClient)

	// HTTP server
	s.Server = kitHttp.NewHttpServer("localhost", "8000")

	s.userAdapter = usersRep.NewAdapter()
	userService := s.userAdapter.GetService()
	userController := users.NewController(authService, userService)

	s.taskAdapter = tasksRep.NewAdapter()
	taskService := s.taskAdapter.GetService()
	taskController := tasks.NewController(taskService)

	s.servAdapter = servRep.NewAdapter()
	deliveryService := s.servAdapter.GetDeliveryService()
	balanceService := s.servAdapter.GetBalanceService()
	servController := services.NewController(balanceService, deliveryService)

	s.bpAdapter = bpRep.NewAdapter()
	bpService := s.bpAdapter.GetService()
	bpController := bp.NewController(bpService)

	s.Server.SetRouters(users.NewRouter(userController), tasks.NewRouter(taskController), services.NewRouter(servController), bp.NewRouter(bpController))

	// session HUB
	s.hub = session.NewHub(s.Server, authService, userService)
	s.Server.SetRouters(s.hub.GetLoginRouteSetter())

	// set auth middlewares
	// the first middleware checks session by X-SESSION-ID header and if correct sets Authorization Bearer with Access Token
	// then the mdw which checks standard Bearer token takes its action
	// TODO: currently if a token expires we don't remove session immediately, but we must
	//s.Server.SetAuthMiddleware(s.hub.SessionMiddleware, s.mdw.CheckToken)

	return s
}

func (s *serviceImpl) Init() error {

	if err := s.userAdapter.Init(); err != nil {
		return err
	}

	if err := s.taskAdapter.Init(); err != nil {
		return err
	}

	if err := s.servAdapter.Init(); err != nil {
		return err
	}

	if err := s.bpAdapter.Init(); err != nil {
		return err
	}

	return nil
}

func (u *serviceImpl) ListenAsync() error {

	go func() {
		log.Fatal(u.Open())
	}()

	return nil
}

func (s *serviceImpl) Close() {
}