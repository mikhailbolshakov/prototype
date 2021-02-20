package config

import (
	"context"
	"encoding/json"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	pb "gitlab.medzdrav.ru/prototype/proto/config"
)

type serviceImpl struct {
	pb.ConfigServiceClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (u *serviceImpl) Get() (*kit.Config, error) {

	l := log.L().Cmp("config").Mth("get")

	rs, err := u.ConfigServiceClient.Get(context.Background(), &pb.ConfigRequest{})
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	var cfg = &kit.Config{}
	if err := json.Unmarshal(rs.Config, cfg); err != nil {
		return nil, err
	}

	if err := kit.EnrichWithEnv("../.env", cfg); err != nil {
		l.E(err).Err("failed to enrich config with local env")
		return nil, err
	}

	return cfg, nil
}
