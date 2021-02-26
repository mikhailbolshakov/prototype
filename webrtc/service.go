package webrtc

import (
	"context"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/domain/impl"
	"gitlab.medzdrav.ru/prototype/webrtc/meta"
	"gitlab.medzdrav.ru/prototype/webrtc/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/webrtc/repository/adapters/ion"
	"gitlab.medzdrav.ru/prototype/webrtc/repository/storage"
)

// NodeId - node id of a service
// TODO: not to hardcode. Should be defined by service discovery procedure
var nodeId = "1"

type serviceImpl struct {
	configAdapter  config.Adapter
	configService  domain.ConfigService
	storageAdapter storage.Adapter
	ionAdapter     ion.Adapter
	ionService     domain.IonService
	webrtcService  domain.WebrtcService
	queue          queue.Queue
}

func New() service.Service {

	s := &serviceImpl{}

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = stan.New()
	s.storageAdapter = storage.NewAdapter()
	strg := s.storageAdapter.GetService()

	s.ionAdapter = ion.NewAdapter()
	s.ionService = s.ionAdapter.GetService()

	s.webrtcService = impl.NewWebrtcService(s.ionService, strg, s.queue)

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

	if err := s.queue.Open(ctx, meta.ServiceCode+nodeId, &queue.Options{
		Url:       c.Nats.Url,
		ClusterId: c.Nats.ClusterId,
	}); err != nil {
		return err
	}

	if err := s.ionAdapter.Init(ctx, c); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync(ctx context.Context) error {
	s.ionAdapter.ListenAsync()
	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {
	s.configAdapter.Close()
	_ = s.queue.Close()
	s.storageAdapter.Close()
	s.ionAdapter.Close(ctx)
}
