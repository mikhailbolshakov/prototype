package config

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/bp/logger"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/kit/log"
	pb "gitlab.medzdrav.ru/prototype/proto/config"
	"os"
	"time"
)

const AWAIT_TIMEOUT = time.Second * 20

type Adapter interface {
	Init(awaitReadiness bool) error
	GetService() domain.ConfigService
	Close()
}

type adapterImpl struct {
	serviceImpl *serviceImpl
	client *kitGrpc.Client
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		serviceImpl: newImpl(),
	}
	return a
}

func (a *adapterImpl) l() log.CLogger {
	return logger.L().Pr("grpc").Cmp("config-adapter")
}

func (a *adapterImpl) Init(awaitReadiness bool) error {

	l := a.l().Mth("init")

	cfgServicePort := os.Getenv("CONFIG_CFG_GRPC_PORT")
	if cfgServicePort == "" {
		err := fmt.Errorf("env var CONFIG_CFG_GRPC_PORT is empty")
		l.E(err).St().Err()
		return err
	}

	cfgServiceHost := os.Getenv("CONFIG_CFG_GRPC_HOST")
	if cfgServiceHost == "" {
		err := fmt.Errorf("env var CONFIG_CFG_GRPC_HOST is empty")
		l.E(err).St().Err()
		return err
	}

	cl, err := kitGrpc.NewClient(cfgServiceHost, cfgServicePort)
	if err != nil {
		l.E(err).St().Err("grpc client error")
		return err
	}

	a.client = cl
	a.serviceImpl.ConfigServiceClient = pb.NewConfigServiceClient(cl.Conn)

	// await remote service starts serving
	if awaitReadiness {
		return a.awaitReadiness()
	}

	return nil

}

func (a *adapterImpl) GetService() domain.ConfigService {
	return a.serviceImpl
}

func (a *adapterImpl) awaitReadiness() error {
	l := a.l().Mth("await").DbgF("awaiting config-server readiness, timeout=%v", AWAIT_TIMEOUT)
	if !a.client.AwaitReadiness(AWAIT_TIMEOUT) {
		err := fmt.Errorf("not ready within timeout")
		l.E(err).Err()
		return err
	}
	l.Dbg("ready")
	return nil
}

func (a *adapterImpl) Close() {
	l := a.l().Mth("close")
	if err := a.client.Conn.Close(); err != nil {
		l.E(err).Err()
	} else {
		l.Dbg("closed")
	}
}