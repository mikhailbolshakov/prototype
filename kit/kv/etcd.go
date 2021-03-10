package kv

import (
	"gitlab.medzdrav.ru/prototype/kit/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"time"
)

type Etcd struct {
	Client *clientv3.Client
	logger log.CLoggerFunc
}

type Options struct {
	Hosts []string
	Dial  *time.Duration
}

func Open(opt *Options, logger log.CLoggerFunc) (*Etcd, error) {

	etcd := &Etcd{
		logger: logger,
	}

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

	logger().Cmp("etcd").Inf("ok")

	etcd.Client = cl
	return etcd, nil

}

func (e *Etcd) Close() error {
	e.logger().Cmp("etcd").Inf("closed")
	if e.Client != nil {
		return e.Client.Close()
	}
	return nil
}
