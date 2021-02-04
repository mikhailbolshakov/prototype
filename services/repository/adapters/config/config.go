package config

import (
	"context"
	"encoding/json"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	pb "gitlab.medzdrav.ru/prototype/proto/config"
	"log"
)

type Service interface {
	Get() (*kit.Config, error)
}

type serviceImpl struct {
	pb.ConfigServiceClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) Get() (*kit.Config, error) {

	rs, err := u.ConfigServiceClient.Get(context.Background(), &pb.ConfigRequest{})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	var cfg = &kit.Config{}
	if err := json.Unmarshal(rs.Config, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
