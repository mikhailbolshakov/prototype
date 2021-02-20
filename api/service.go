package api

import (
	"context"
	"github.com/Nerzal/gocloak/v7"
	"gitlab.medzdrav.ru/prototype/api/public"
	"gitlab.medzdrav.ru/prototype/api/public/bp"
	"gitlab.medzdrav.ru/prototype/api/public/chat"
	"gitlab.medzdrav.ru/prototype/api/public/monitoring"
	"gitlab.medzdrav.ru/prototype/api/public/services"
	"gitlab.medzdrav.ru/prototype/api/public/tasks"
	"gitlab.medzdrav.ru/prototype/api/public/users"
	bpRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/bp"
	chatRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/api/repository/adapters/config"
	servRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/services"
	tasksRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/tasks"
	usersRep "gitlab.medzdrav.ru/prototype/api/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/api/session"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/proto"
)

// NodeId - node id of a service
// TODO: not to hardcode. Should be defined by service discovery procedure
var nodeId = "1"
const serviceCode = "api"

type serviceImpl struct {
	http            *kitHttp.Server
	keycloak        gocloak.GoCloak
	authMdw         auth.Middleware
	hub             session.Hub
	userAdapter     usersRep.Adapter
	userService     public.UserService
	taskAdapter     tasksRep.Adapter
	taskService     public.TaskService
	chatAdapter     chatRep.Adapter
	chatService     public.ChatService
	servAdapter     servRep.Adapter
	balanceService  public.BalanceService
	deliveryService public.DeliveryService
	bpAdapter       bpRep.Adapter
	bpService       public.BpService
	configAdapter   config.Adapter
	configService   public.ConfigService
	queue           queue.Queue
	queueListener   listener.QueueListener
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

	s.chatAdapter = chatRep.NewAdapter(s.userService)
	s.chatService = s.chatAdapter.GetService()

	s.queue = stan.New()
	s.queueListener = listener.NewQueueListener(s.queue)

	return s
}

func (s *serviceImpl) initHttpServer(ctx context.Context, c *kitConfig.Config) error {

	authClient := &auth.ClientSecurityInfo{
		ID:     c.Keycloak.ClientId,
		Secret: c.Keycloak.ClientSecret,
		Realm:  c.Keycloak.ClientRealm,
	}

	s.keycloak = gocloak.NewClient(c.Keycloak.Url)

	s.authMdw = auth.NewMdw(ctx, s.keycloak, authClient, "", "")

	authService := auth.New(ctx, s.keycloak, authClient)

	s.http = kitHttp.NewHttpServer(c.Http.Host, c.Http.Port, c.Http.Tls.Cert, c.Http.Tls.Key)

	// session HUB
	s.hub = session.NewHub(c, s.http, authService, s.userService, s.chatService)
	// setup a NATS message handler on events forwarded to websocket
	s.queueListener.Add(queue.QUEUE_TYPE_AT_MOST_ONCE, proto.QUEUE_TOPIC_OUTGOING_WS_EVENT, s.hub.GetOutgoingWsEventsHandler())

	s.http.SetRouters(s.hub.GetLoginRouteSetter())

	// set auth middlewares
	// the first middleware checks session by X-SESSION-ID header and if correct sets Authorization Bearer with Access Token
	// then the mdw which checks standard Bearer token takes its action
	// TODO: currently if a token expires we don't remove session immediately, but we must
	s.http.SetAuthMiddleware(s.hub.SessionMiddleware, s.authMdw.CheckToken)
	// middleware for routes that don't require auth and don't have sessions like login (hm... we have to think it over)
	s.http.SetNoAuthMiddleware(s.hub.NoSessionMiddleware)

	userController := users.NewController(authService, s.userService)
	taskController := tasks.NewController(s.taskService)
	servController := services.NewController(s.balanceService, s.deliveryService)
	bpController := bp.NewController(s.bpService)
	chatController := chat.NewController(s.chatService)
	monitorController := monitoring.NewController(s.hub.GetMonitor())

	s.http.SetRouters(users.NewRouter(userController),
		tasks.NewRouter(taskController),
		services.NewRouter(servController),
		bp.NewRouter(bpController),
		chat.NewRouter(chatController),
		monitoring.NewRouter(monitorController))

	return nil
}

func (s *serviceImpl) Init(ctx context.Context) error {

	if err := s.configAdapter.Init(); err != nil {
		return err
	}

	c, err := s.configService.Get()
	if err != nil {
		return err
	}

	if err := s.initHttpServer(ctx, c); err != nil {
		return err
	}

	if err := s.queue.Open(ctx, serviceCode + nodeId, &queue.Options{
		Url:       c.Nats.Url,
		ClusterId: c.Nats.ClusterId,
	}); err != nil {
		return err
	}

	if err := s.userAdapter.Init(c); err != nil {
		return err
	}

	if err := s.taskAdapter.Init(c); err != nil {
		return err
	}

	if err := s.servAdapter.Init(c); err != nil {
		return err
	}

	if err := s.chatAdapter.Init(c); err != nil {
		return err
	}

	if err := s.bpAdapter.Init(c); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {

	s.http.Listen()
	s.queueListener.ListenAsync()

	return nil
}

func (s *serviceImpl) Close(context.Context) {
	s.bpAdapter.Close()
	s.servAdapter.Close()
	s.userAdapter.Close()
	s.taskAdapter.Close()
	s.chatAdapter.Close()
	s.configAdapter.Close()
	s.http.Close()
}
