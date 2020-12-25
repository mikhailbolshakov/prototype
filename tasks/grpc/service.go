package grpc

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	domain "gitlab.medzdrav.ru/prototype/tasks/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	domain domain.TaskService
	pb.UnimplementedTasksServer
}

func New() *Service {
	s := &Service{
		domain: domain.NewTaskService(),
	}
	return s
}

func (s *Service) New(ctx context.Context, rq *pb.NewTaskRequest) (*pb.Task, error) {

	task, err := s.fromPb(rq)
	if err != nil {
		return nil, err
	}

	task, err = s.domain.New(task)
	if err != nil {
		return nil, err
	}

	rs, err := s.fromDomain(task)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (s *Service) NextTransitions(ctx context.Context, rq *pb.NextTransitionsRequest) (*pb.NextTransitionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NextTransitions not implemented")
}
func (s *Service) MakeTransition(ctx context.Context, rq *pb.MakeTransitionRequest) (*pb.Task, error) {

	task, err := s.domain.MakeTransition(rq.TaskId, rq.TransitionId)
	if err != nil {
		return nil, err
	}

	rs, err := s.fromDomain(task)
	if err != nil {
		return nil, err
	}

	return rs, nil
}