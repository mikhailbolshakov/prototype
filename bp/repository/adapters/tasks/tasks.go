package tasks

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/queue_model"
	"log"
)

type Service interface {
	GetByChannelId(channelId string) []*pb.Task
	CreateTask(rq *pb.NewTaskRequest) (*pb.Task, error)
	MakeTransition(rq *pb.MakeTransitionRequest) error
}

type TaskHandler func(task *queue_model.Task)

type serviceImpl struct {
	pb.TasksClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{
	}
	return a
}

func (u *serviceImpl) GetByChannelId(channelId string) []*pb.Task {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs, err := u.GetByChannel(ctx, &pb.GetByChannelRequest{ChannelId: channelId})
	if err != nil {
		log.Printf("error: %v", err)
		return []*pb.Task{}
	}
	return rs.Tasks
}

func (u *serviceImpl) MakeTransition(rq *pb.MakeTransitionRequest) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := u.TasksClient.MakeTransition(ctx, rq)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	return nil
}

func (u *serviceImpl) CreateTask(rq *pb.NewTaskRequest) (*pb.Task, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return u.New(ctx, rq)
}

