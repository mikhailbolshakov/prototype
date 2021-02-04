package tasks

import (
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
	GetService() Service
	Close()
}

type adapterImpl struct {
	taskServiceImpl *serviceImpl
	client          *kitGrpc.Client
}

func NewAdapter(queue queue.Queue) Adapter {
	a := &adapterImpl{
		taskServiceImpl: newImpl(),
	}
	return a
}

func (a *adapterImpl) Init(c *kitConfig.Config) error {
	cfg := c.Services["tasks"]
	cl, err := kitGrpc.NewClient(cfg.Grpc.Hosts[0], cfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.taskServiceImpl.TasksClient = pb.NewTasksClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() Service {
	return a.taskServiceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}
