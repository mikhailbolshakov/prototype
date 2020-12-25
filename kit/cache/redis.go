package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
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
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", params.Host, params.Port),
		Password: params.Password,
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to Redis")
	return &Redis{
		Instance: client,
		Ttl:      time.Duration(params.Ttl) * time.Second,
	}, nil
}
