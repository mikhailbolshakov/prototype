package config

import (
	"context"
	"encoding/json"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	pb "gitlab.medzdrav.ru/prototype/proto/config"
)

type serviceImpl struct {
	pb.ConfigServiceClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) Get(ctx context.Context) (*kit.Config, error) {

	rs, err := u.ConfigServiceClient.Get(ctx, &pb.ConfigRequest{})
	if err != nil {
		return nil, err
	}

	var cfg = &kit.Config{}
	if err := json.Unmarshal(rs.Config, cfg); err != nil {
		return nil, err
	}

	if err := kit.EnrichWithEnv("../.env", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
