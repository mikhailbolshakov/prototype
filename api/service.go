package api

import (
	"context"
	"github.com/Nerzal/gocloak/v7"
	"gitlab.medzdrav.ru/prototype/api/public/bp"
	"gitlab.medzdrav.ru/prototype/api/public/services"
	"gitlab.medzdrav.ru/prototype/api/public/tasks"
	"gitlab.medzdrav.ru/prototype/api/public/users"
	bpRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/bp"
	"gitlab.medzdrav.ru/prototype/api/repository/adapters/config"
	servRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/services"
	tasksRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/tasks"
	usersRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/api/session"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"log"
)

type serviceImpl struct {
	*kitHttp.Server
	keycloak        gocloak.GoCloak
	mdw             auth.Middleware
	hub             session.Hub
	userAdapter     usersRep.Adapter
	userService     usersRep.Service
	taskAdapter     tasksRep.Adapter
	taskService     tasksRep.Service
	servAdapter     servRep.Adapter
	balanceService  servRep.BalanceService
	deliveryService servRep.DeliveryService
	bpAdapter       bpRep.Adapter
	bpService       bpRep.Service
	configAdapter   config.Adapter
	configService   config.Service
}

func New() service.Service {
	s := &serviceImpl{}

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.userAdapter = usersRep.NewAdapter()
	s.userService = s.userAdapter.GetService()

	s.taskAdapter = tasksRep.NewAdapter()
	s.taskService = s.taskAdapter.GetService()

	s.servAdapter = servRep.NewAdapter()
	s.deliveryService = s.servAdapter.GetDeliveryService()
	s.balanceService = s.servAdapter.GetBalanceService()

	s.bpAdapter = bpRep.NewAdapter()
	s.bpService = s.bpAdapter.GetService()

	return s
}

func (s *serviceImpl) Init() error {

	ctx := context.Background()

	if err := s.configAdapter.Init(); err != nil {
		return err
	}

	c, err := s.configService.Get()
	if err != nil {
		return err
	}

	authClient := &auth.ClientSecurityInfo{
		ID:     c.Keycloak.ClientId,
		Secret: c.Keycloak.ClientSecret,
		Realm:  c.Keycloak.ClientRealm,
	}

	s.keycloak = gocloak.NewClient(c.Keycloak.Url)
	s.mdw = auth.NewMdw(ctx, s.keycloak, authClient, "", "")

	authService := auth.New(ctx, s.keycloak, authClient)

	// HTTP server
	s.Server = kitHttp.NewHttpServer(c.Http.Host, c.Http.Port)

	userController := users.NewController(authService, s.userService)
	taskController := tasks.NewController(s.taskService)
	servController := services.NewController(s.balanceService, s.deliveryService)
	bpController := bp.NewController(s.bpService)

	s.Server.SetRouters(users.NewRouter(userController), tasks.NewRouter(taskController), services.NewRouter(servController), bp.NewRouter(bpController))

	// session HUB
	s.hub = session.NewHub(c, s.Server, authService, s.userService)
	s.Server.SetRouters(s.hub.GetLoginRouteSetter())

	// set auth middlewares
	// the first middleware checks session by X-SESSION-ID header and if correct sets Authorization Bearer with Access Token
	// then the mdw which checks standard Bearer token takes its action
	// TODO: currently if a token expires we don't remove session immediately, but we must
	//s.Server.SetAuthMiddleware(s.hub.SessionMiddleware, s.mdw.CheckToken)

	if err := s.userAdapter.Init(c); err != nil {
		return err
	}

	if err := s.taskAdapter.Init(c); err != nil {
		return err
	}

	if err := s.servAdapter.Init(c); err != nil {
		return err
	}

	if err := s.bpAdapter.Init(c); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) ListenAsync() error {

	go func() {
		log.Fatal(s.Open())
	}()

	return nil
}

func (s *serviceImpl) Close() {
	s.bpAdapter.Close()
	s.servAdapter.Close()
	s.userAdapter.Close()
	s.taskAdapter.Close()
	s.configAdapter.Close()
	_ = s.Srv.Close()
}
