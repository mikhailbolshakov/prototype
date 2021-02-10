package config

import (
	"context"
	"encoding/json"
	"fmt"
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

	rs, err := u.ConfigServiceClient.Get(context.Background(), &pb.ConfigRequest{})
	if err != nil {
		log.Err(err, true)
		return nil, err
	}

	var cfg = &kit.Config{}
	if err := json.Unmarshal(rs.Config, cfg); err != nil {
		return nil, err
	}

	if err := kit.EnrichWithEnv("../.env", cfg); err != nil {
		log.Err(fmt.Errorf("failed to enrich config with local env"), false)
		return nil, err
	}

	return cfg, nil
}
