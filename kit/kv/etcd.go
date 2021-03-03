package kv

import (
	"gitlab.medzdrav.ru/prototype/kit/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"time"
)

type Etcd struct {
	Client *clientv3.Client
}

type Options struct {
	Hosts []string
	Dial  *time.Duration
}

func Open(opt *Options) (*Etcd, error) {

	etcd := &Etcd{}

	var dial = time.Second * 3

	if opt.Dial != nil {
		dial = *opt.Dial
	}

	cl, err := clientv3.New(clientv3.Config{
		DialTimeout: dial,
		DialOptions: []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()},
		Endpoints:   opt.Hosts,
	})
	if err != nil {
		return nil, err
	}

	log.L().Pr("etcd").Inf("ok")

	etcd.Client = cl
	return etcd, nil

}

func (e *Etcd) Close() error {
	log.L().Pr("etcd").Inf("closed")
	if e.Client != nil {
		return e.Client.Close()
	}
	return nil
}
