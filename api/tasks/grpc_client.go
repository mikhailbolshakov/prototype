package tasks

import (
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
)

type grpcClient struct {
	*kitGrpc.Client
	tasks pb.TasksClient
}

func newGrpcClient() (*grpcClient, error) {

	c := &grpcClient{}
	cl, err := kitGrpc.NewClient("localhost", "50052")
	if err != nil {
		return nil, err
	}
	c.Client = cl
	c.tasks = pb.NewTasksClient(c.Conn)

	return c, nil

}
