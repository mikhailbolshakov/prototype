package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"time"
)

type Redis struct {
	Instance *redis.Client
	Ttl      time.Duration
}

type Params struct {
	Host     string
	Port     string
	Password string
	Ttl      uint
}

func Open(params *Params) (*Redis, error) {

	l := log.L().Cmp("redis").Mth("open")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", params.Host, params.Port),
		Password: params.Password,
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	l.Inf("ok")
	return &Redis{
		Instance: client,
		Ttl:      time.Duration(params.Ttl) * time.Second,
	}, nil
}

func (r *Redis) Close() {
	if r.Instance != nil {
		_ = r.Instance.Close()
	}
}