package users

import (
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
)

type grpcClient struct {
	*kitGrpc.Client
	users pb.UsersClient
}

func newGrpcClient() (*grpcClient, error) {

	c := &grpcClient{}
	cl, err := kitGrpc.NewClient("localhost", "50051")
	if err != nil {
		return nil, err
	}
	c.Client = cl
	c.users = pb.NewUsersClient(c.Conn)

	return c, nil

}
