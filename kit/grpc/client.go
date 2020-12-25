package grpc

import (
	"fmt"
	"google.golang.org/grpc"
)

type Client struct {
	Conn *grpc.ClientConn
}

func NewClient(host, port string) (*Client, error) {

	gc, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure() /*, grpc.WithBlock()*/)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Conn: gc,
	}

	return c, nil
}