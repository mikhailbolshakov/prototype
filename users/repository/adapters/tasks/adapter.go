package tasks

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
)

type Adapter interface {
	Init() error
	GetService() Service
}

type adapterImpl struct {
	taskServiceImpl *serviceImpl
	initialized bool
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		taskServiceImpl: newImpl(),
		initialized: false,
	}
	return a
}

func (a *adapterImpl) Init() error {
	cl, err := kitGrpc.NewClient("localhost", "50052")
	if err != nil {
		return err
	}
	a.taskServiceImpl.TasksClient = pb.NewTasksClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() Service {
	return a.taskServiceImpl
}
