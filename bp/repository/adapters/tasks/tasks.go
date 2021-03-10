package tasks

import (
	"context"
	"gitlab.medzdrav.ru/prototype/bp/logger"
	"gitlab.medzdrav.ru/prototype/kit/log"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
)

type serviceImpl struct {
	pb.TasksClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{
	}
	return a
}

func (u *serviceImpl) l() log.CLogger {
	return logger.L().Cmp("task-adapter")
}

func (u *serviceImpl) GetByChannelId(ctx context.Context, channelId string) []*pb.Task {

	l := u.l().Mth("get-by-channel").C(ctx).F(log.FF{"channel": channelId}).Dbg()

	rs, err := u.GetByChannel(ctx, &pb.GetByChannelRequest{ChannelId: channelId})
	if err != nil {
		l.E(err).Err()
		return []*pb.Task{}
	}
	return rs.Tasks
}

func (u *serviceImpl) MakeTransition(ctx context.Context, rq *pb.MakeTransitionRequest) error {

	l := u.l().Mth("make-transition").C(ctx).F(log.FF{"task": rq.TaskId, "tr": rq.TransitionId}).Dbg()

	_, err := u.TasksClient.MakeTransition(ctx, rq)
	if err != nil {
		l.E(err).Err()
		return err
	}
	return nil
}

func (u *serviceImpl) New(ctx context.Context, rq *pb.NewTaskRequest) (*pb.Task, error) {
	u.l().Mth("new").C(ctx).Dbg()
	return u.TasksClient.New(ctx, rq)
}

func (u *serviceImpl) Search(ctx context.Context, rq *pb.SearchRequest) ([]*pb.Task, error) {

	u.l().Mth("search").C(ctx).Dbg()

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

