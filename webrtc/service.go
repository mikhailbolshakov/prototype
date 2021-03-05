package webrtc

import (
	"context"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/domain/impl"
	"gitlab.medzdrav.ru/prototype/webrtc/grpc"
	"gitlab.medzdrav.ru/prototype/webrtc/jsonrpc"
	"gitlab.medzdrav.ru/prototype/webrtc/meta"
	"gitlab.medzdrav.ru/prototype/webrtc/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/webrtc/repository/adapters/sessions"
	"gitlab.medzdrav.ru/prototype/webrtc/repository/storage"
)

type serviceImpl struct {
	http            *kitHttp.Server
	grpc            *grpc.Server
	jsonRpc         *jsonrpc.Server
	configAdapter   config.Adapter
	configService   domain.ConfigService
	storageAdapter  storage.Adapter
	webrtcService   domain.WebrtcService
	queue           queue.Queue
	sessionsAdapter sessions.Adapter
}

func New() service.Service {

	s := &serviceImpl{}

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = stan.New()
	s.storageAdapter = storage.NewAdapter()

	s.sessionsAdapter = sessions.NewAdapter()

	s.webrtcService = impl.NewWebrtcService(s.storageAdapter.GetRoomCoordinator(), s.storageAdapter.GetService(), s.queue)

	return s
}

func (s *serviceImpl) Init(ctx context.Context) error {

	if err := s.configAdapter.Init(); err != nil {
		return err
	}

	c, err := s.configService.Get(ctx)
	if err != nil {
		return err
	}

	if err := s.storageAdapter.Init(c); err != nil {
		return err
	}

	if err := s.sessionsAdapter.Init(c); err != nil {
		return err
	}

	if err := s.queue.Open(ctx, meta.ServiceCode+meta.NodeId, &queue.Options{
		Url:       c.Nats.Url,
		ClusterId: c.Nats.ClusterId,
	}); err != nil {
		return err
	}

	if err := s.webrtcService.Init(ctx, c); err != nil {
		return err
	}

	s.grpc = grpc.New(s.webrtcService)
	if err := s.grpc.Init(c); err != nil {
		return err
	}

	s.http = kitHttp.NewHttpServer(c.Webrtc.Signal.Host, c.Webrtc.Signal.Port, "", "")

	s.jsonRpc = jsonrpc.New(s.http, s.sessionsAdapter.GetService(), s.webrtcService)

	return nil

}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {
	s.grpc.ListenAsync()
	s.http.Listen()
	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {
	s.grpc.Close()
	s.http.Close()
	s.configAdapter.Close()
	_ = s.queue.Close()
	s.storageAdapter.Close()
	s.sessionsAdapter.Close()
}
