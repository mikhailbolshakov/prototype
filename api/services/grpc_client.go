package services

import (
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
)

type grpcClient struct {
	*kitGrpc.Client
	delivery pb.DeliveryServiceClient
	balance pb.BalanceServiceClient
}

func newGrpcClient() (*grpcClient, error) {

	c := &grpcClient{}
	cl, err := kitGrpc.NewClient("localhost", "50054")
	if err != nil {
		return nil, err
	}
	c.Client = cl
	c.delivery = pb.NewDeliveryServiceClient(c.Conn)
	c.balance = pb.NewBalanceServiceClient(c.Conn)

	return c, nil

}
