package config

import (
	"errors"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/config"
	"gitlab.medzdrav.ru/prototype/services/domain"
)

type Adapter interface {
	Init() error
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

func (a *adapterImpl) Init() error {

	if envs, err := kitConfig.Env("../.env"); err == nil {

		port, ok := envs["CONFIG_CFG_GRPC_PORT"]
		if !ok {
			return errors.New("config server port isn't specified")
		}
		host, ok := envs["CONFIG_CFG_GRPC_HOST"]
		if !ok {
			return errors.New("config server port isn't specified")
		}

		cl, err := kitGrpc.NewClient(host, port)
		if err != nil {
			return err
		}

		a.client = cl
		a.serviceImpl.ConfigServiceClient = pb.NewConfigServiceClient(cl.Conn)

		return nil

	} else {
		return err
	}

}

func (a *adapterImpl) GetService() domain.ConfigService {
	return a.serviceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}