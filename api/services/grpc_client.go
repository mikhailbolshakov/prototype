package services

import (
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
)

type grpcClient struct {
	*kitGrpc.Client
	services pb.UserServicesClient
}

func newGrpcClient() (*grpcClient, error) {

	c := &grpcClient{}
	cl, err := kitGrpc.NewClient("localhost", "50054")
	if err != nil {
		return nil, err
	}
	c.Client = cl
	c.services = pb.NewUserServicesClient(c.Conn)

	return c, nil

}
