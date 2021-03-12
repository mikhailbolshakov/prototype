package webrtc

import (
	"context"
	"fmt"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/domain/impl"
	"gitlab.medzdrav.ru/prototype/webrtc/grpc"
	"gitlab.medzdrav.ru/prototype/webrtc/jsonrpc"
	"gitlab.medzdrav.ru/prototype/webrtc/logger"
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
	roomService     domain.RoomService
	storageAdapter  storage.Adapter
	webrtcService   domain.WebrtcService
	queue           queue.Queue
	sessionsAdapter sessions.Adapter
}

func New() service.Service {

	s := &serviceImpl{}

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = stan.New(logger.LF())
	s.storageAdapter = storage.NewAdapter()

	s.roomService = impl.NewRoomService(s.storageAdapter.GetService())
	s.sessionsAdapter = sessions.NewAdapter()

	s.webrtcService = impl.NewWebrtcService(s.storageAdapter.GetRoomCoordinator(), s.roomService, s.queue)

	return s
}

func (s *serviceImpl) GetCode() string {
	return meta.Meta.ServiceCode()
}

func (s *serviceImpl) Init(ctx context.Context) error {

	if err := s.configAdapter.Init(true); err != nil {
		return err
	}

	cfg, err := s.configService.Get(ctx)
	if err != nil {
		return err
	}

	if srvCfg, ok := cfg.Services[meta.Meta.ServiceCode()]; ok {
		logger.Logger.SetLevel(srvCfg.Log.Level)
	} else {
		return fmt.Errorf("service config isn't specified")
	}

	if err := s.storageAdapter.Init(cfg); err != nil {
		return err
	}

	if err := s.sessionsAdapter.Init(cfg); err != nil {
		return err
	}

	if err := s.queue.Open(ctx, meta.Meta.InstanceId(), &queue.Options{
		Url:       cfg.Nats.Url,
		ClusterId: cfg.Nats.ClusterId,
	}); err != nil {
		return err
	}

	if err := s.webrtcService.Init(ctx, cfg); err != nil {
		return err
	}

	s.grpc = grpc.New(s.webrtcService, s.roomService)
	if err := s.grpc.Init(cfg); err != nil {
		return err
	}

	s.http = kitHttp.NewHttpServer(cfg.Webrtc.Signal.Host, cfg.Webrtc.Signal.Port, "", "", logger.LF())

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
