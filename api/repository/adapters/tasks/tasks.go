package tasks

import (
	"context"
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

func (u *serviceImpl) MakeTransition(ctx context.Context, rq *pb.MakeTransitionRequest) (*pb.Task, error) {
	return u.TasksClient.MakeTransition(ctx, rq)
}

func (u *serviceImpl) New(ctx context.Context, rq *pb.NewTaskRequest) (*pb.Task, error) {
	return u.TasksClient.New(ctx, rq)
}

func (u *serviceImpl) SetAssignee(ctx context.Context, rq *pb.SetAssigneeRequest) (*pb.Task, error) {
	return u.TasksClient.SetAssignee(ctx, rq)
}

func (u *serviceImpl) GetById(ctx context.Context, id string) (*pb.Task, error) {
	return u.TasksClient.GetById(ctx, &pb.GetByIdRequest{Id: id})
}

func (u *serviceImpl) Search(ctx context.Context, rq *pb.SearchRequest) (*pb.SearchResponse, error) {
	return u.TasksClient.Search(ctx, rq)
}

func (u *serviceImpl) GetAssignmentLog(ctx context.Context, rq *pb.AssignmentLogRequest) (*pb.AssignmentLogResponse, error) {
	return u.TasksClient.GetAssignmentLog(ctx, rq)
}

func (u *serviceImpl) GetHistory(ctx context.Context, taskId string) (*pb.GetHistoryResponse, error) {
	return u.TasksClient.GetHistory(ctx, &pb.GetHistoryRequest{TaskId: taskId})
}

