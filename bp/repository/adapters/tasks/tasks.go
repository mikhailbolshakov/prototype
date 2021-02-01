package tasks

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/queue_model"
	"log"
)

type Service interface {
	GetByChannelId(channelId string) []*pb.Task
	New(rq *pb.NewTaskRequest) (*pb.Task, error)
	MakeTransition(rq *pb.MakeTransitionRequest) error
	Search(rq *pb.SearchRequest) ([]*pb.Task, error)
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

func (u *serviceImpl) New(rq *pb.NewTaskRequest) (*pb.Task, error) {
	return u.TasksClient.New(context.Background(), rq)
}

func (u *serviceImpl) Search(rq *pb.SearchRequest) ([]*pb.Task, error) {

	if rq.Status == nil {
		rq.Status = &pb.Status{}
	}

	if rq.Type == nil {
		rq.Type = &pb.Type{}
	}

	if rq.Assignee == nil {
		rq.Assignee = &pb.Assignee{}
	}

	if rq.Paging == nil {
		rq.Paging = &pb.PagingRequest{}
	}

	if rs, err := u.TasksClient.Search(context.Background(), rq); err != nil {
		return nil, err
	} else {
		return rs.Tasks, nil
	}
}

