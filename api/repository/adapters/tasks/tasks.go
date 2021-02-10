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

func (u *serviceImpl) MakeTransition(rq *pb.MakeTransitionRequest) (*pb.Task, error) {
	return u.TasksClient.MakeTransition(context.Background(), rq)
}

func (u *serviceImpl) New(rq *pb.NewTaskRequest) (*pb.Task, error) {
	return u.TasksClient.New(context.Background(), rq)
}

func (u *serviceImpl) SetAssignee(rq *pb.SetAssigneeRequest) (*pb.Task, error) {
	return u.TasksClient.SetAssignee(context.Background(), rq)
}

func (u *serviceImpl) GetById(id string) (*pb.Task, error) {
	return u.TasksClient.GetById(context.Background(), &pb.GetByIdRequest{Id: id})
}

func (u *serviceImpl) Search(rq *pb.SearchRequest) (*pb.SearchResponse, error) {
	return u.TasksClient.Search(context.Background(), rq)
}

func (u *serviceImpl) GetAssignmentLog(rq *pb.AssignmentLogRequest) (*pb.AssignmentLogResponse, error) {
	return u.TasksClient.GetAssignmentLog(context.Background(), rq)
}

func (u *serviceImpl) GetHistory(taskId string) (*pb.GetHistoryResponse, error) {
	return u.TasksClient.GetHistory(context.Background(), &pb.GetHistoryRequest{TaskId: taskId})
}

