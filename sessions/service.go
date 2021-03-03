package sessions

import (
	"context"
	"github.com/Nerzal/gocloak/v7"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/proto"
	"gitlab.medzdrav.ru/prototype/sessions/domain"
	"gitlab.medzdrav.ru/prototype/sessions/domain/impl"
	"gitlab.medzdrav.ru/prototype/sessions/grpc"
	"gitlab.medzdrav.ru/prototype/sessions/meta"
	"gitlab.medzdrav.ru/prototype/sessions/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/sessions/repository/adapters/config"
	metrics "gitlab.medzdrav.ru/prototype/sessions/repository/adapters/metrics"
	"gitlab.medzdrav.ru/prototype/sessions/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/sessions/repository/storage"
)

// NodeId - node id of a service
// TODO: not to hardcode. Should be defined by service discovery procedure
var nodeId = "1"

type serviceImpl struct {
	http            *kitHttp.Server
	keycloak        gocloak.GoCloak
	authMdw         auth.Middleware
	sessionsService domain.SessionsService
	monitorService  domain.SessionMonitor
	configAdapter   config.Adapter
	cfgService      domain.CfgService
	grpc            *grpc.Server
	usersAdapter    users.Adapter
	chatAdapter     chat.Adapter
	queue           queue.Queue
	queueListener   listener.QueueListener
	storageAdapter  storage.Adapter
	metricsAdapter  metrics.Adapter
}

func New() service.Service {

	s := &serviceImpl{}

	s.queue = stan.New()
	s.queueListener = listener.NewQueueListener(s.queue)

	s.configAdapter = config.NewAdapter()
	s.cfgService = s.configAdapter.GetService()
	s.usersAdapter = users.NewAdapter()
	s.storageAdapter = storage.NewAdapter()
	s.metricsAdapter = metrics.NewAdapter()
	s.chatAdapter = chat.NewAdapter()

	return s
}

func (s *serviceImpl) Init(ctx context.Context) error {

	if err := s.configAdapter.Init(); err != nil {
		return err
	}

	c, err := s.cfgService.Get(ctx)
	if err != nil {
		return err
	}

	authClient := &auth.ClientSecurityInfo{
		ID:     c.Keycloak.ClientId,
		Secret: c.Keycloak.ClientSecret,
		Realm:  c.Keycloak.ClientRealm,
	}

	s.keycloak = gocloak.NewClient(c.Keycloak.Url)

	s.authMdw = auth.NewMdw(ctx, s.keycloak, authClient, "", "")

	authService := auth.New(ctx, s.keycloak, authClient)

	s.http = kitHttp.NewHttpServer(c.Http.Host, c.Http.WsPort, c.Http.Tls.Cert, c.Http.Tls.Key)

	// session HUB
	hub := impl.NewHub(c, s.http, s.metricsAdapter.GetService())
	s.sessionsService = impl.NewSessionsService(hub, authService, s.usersAdapter.GetUserService(), s.chatAdapter.GetService())
	s.monitorService = impl.NewMonitorService(hub, s.metricsAdapter.GetService())

	// setup a NATS message handler on events forwarded to websocket
	s.queueListener.Add(queue.QUEUE_TYPE_AT_MOST_ONCE, proto.QUEUE_TOPIC_OUTGOING_WS_EVENT, hub.GetOutgoingWsEventsHandler())

	s.grpc = grpc.New(s.sessionsService, s.monitorService)

	if err := s.storageAdapter.Init(c); err != nil {
		return err
	}

	if err := s.grpc.Init(c); err != nil {
		return err
	}

	if err := s.usersAdapter.Init(c); err != nil {
		return err
	}

	if err := s.chatAdapter.Init(c); err != nil {
		return err
	}

	if err := s.metricsAdapter.Init(c); err != nil {
		return err
	}

	if err := s.queue.Open(ctx, meta.ServiceCode+nodeId, &queue.Options{
		Url:       c.Nats.Url,
		ClusterId: c.Nats.ClusterId,
	}); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {

	s.http.Listen()
	s.grpc.ListenAsync()
	s.queueListener.ListenAsync()

	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {

	s.configAdapter.Close()

	s.metricsAdapter.Close()
	s.chatAdapter.Close()
	s.usersAdapter.Close()

	_ = s.queue.Close()
	s.storageAdapter.Close()
	s.grpc.Close()
}
