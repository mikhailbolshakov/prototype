package config

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/api/logger"
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

func (u *serviceImpl) Get() (*pb.Config, error) {

	l := logger.L().Cmp("config").Mth("get")

	rs, err := u.ConfigServiceClient.Get(context.Background(), &pb.ConfigRequest{})
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	var cfg = &pb.Config{}
	if err := json.Unmarshal(rs.Config, cfg); err != nil {
		return nil, err
	}

	if err := kit.EnrichWithEnv("../.env", cfg); err != nil {
		l.E(err).Err("failed to enrich config with local env")
		return nil, err
	}

	return cfg, nil
}
