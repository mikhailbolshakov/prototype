package ion

import (
	"context"
	"github.com/pion/ion-sfu/pkg/sfu"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
)

type Adapter interface {
	Init(ctx context.Context, cfg *kitConfig.Config) error
	GetService() domain.IonService
	Close(ctx context.Context)
	ListenAsync()
}

type adapterImpl struct {
	ion *ionImpl
	signalErrChan chan error
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		ion: &ionImpl{},
	}
	return a
}

func (a *adapterImpl) listenSignalErr() {
	go func() {
		l := log.L().Cmp("webrtc").Mth("signal-err")
		for err := range a.signalErrChan {
			l.E(err).Err()
		}
	}()
}

func (a *adapterImpl) Init(ctx context.Context, cfg *kitConfig.Config) error {

	a.ion.sfu = sfu.NewSFU(*cfg.Webrtc.SFU)

	if c, err := newCoordinator(ctx, cfg); err == nil {
		a.ion.coordinator = c
	} else {
		return err
	}

	a.ion.signal, a.signalErrChan = newSignal(a.ion.sfu, a.ion.coordinator, cfg.Webrtc.Signal)

	return nil
}

func (a *adapterImpl) ListenAsync() {
	go a.ion.signal.ServeWebsocket()
}

func (a *adapterImpl) GetService() domain.IonService {
	return a.ion
}

func (a *adapterImpl) Close(ctx context.Context) {
	close(a.signalErrChan)
}
