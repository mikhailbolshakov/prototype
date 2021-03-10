package grpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type Client struct {
	*readinessAwaiter
	Conn *grpc.ClientConn
}

func NewClient(host, port string) (*Client, error) {

	gc, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(grpc_middleware.ChainUnaryClient(ContextUnaryClientInterceptor())))
		grpc.WithChainStreamInterceptor(grpc_middleware.ChainStreamClient(ContextStreamClientInterceptor()))
	if err != nil {
		return nil, err
	}

	c := &Client{
		Conn: gc,
	}

	c.readinessAwaiter = newReadinessAwaiter(gc)

	return c, nil
}