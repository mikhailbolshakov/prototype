package tasks

import (
	"gitlab.medzdrav.ru/prototype/bp/domain"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
)

type Adapter interface {
	Init(c *config.Config) error
	GetService() domain.TaskService
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

func (a *adapterImpl) Init(c *config.Config) error {
	cfg := c.Services["tasks"]
	cl, err := kitGrpc.NewClient(cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		return err
	}
	a.client = cl
	a.taskServiceImpl.TasksClient = pb.NewTasksClient(cl.Conn)
	return nil
}

func (a *adapterImpl) GetService() domain.TaskService {
	return a.taskServiceImpl
}

func (a *adapterImpl) Close() {
	_ = a.client.Conn.Close()
}
