package tasks

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"log"
)

type serviceImpl struct {
	pb.TasksClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{
	}
	return a
}

func (u *serviceImpl) GetByChannelId(ctx context.Context, channelId string) []*pb.Task {
	rs, err := u.GetByChannel(ctx, &pb.GetByChannelRequest{ChannelId: channelId})
	if err != nil {
		log.Printf("error: %v", err)
		return []*pb.Task{}
	}
	return rs.Tasks
}

func (u *serviceImpl) MakeTransition(ctx context.Context, rq *pb.MakeTransitionRequest) error {
	_, err := u.TasksClient.MakeTransition(ctx, rq)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	return nil
}

func (u *serviceImpl) New(ctx context.Context, rq *pb.NewTaskRequest) (*pb.Task, error) {
	return u.TasksClient.New(ctx, rq)
}

func (u *serviceImpl) Search(ctx context.Context, rq *pb.SearchRequest) ([]*pb.Task, error) {

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

	if rs, err := u.TasksClient.Search(ctx, rq); err != nil {
		return nil, err
	} else {
		return rs.Tasks, nil
	}
}

