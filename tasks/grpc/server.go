package grpc

import (
	"context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	domain "gitlab.medzdrav.ru/prototype/tasks/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	*kitGrpc.Server
	domain domain.TaskService
	search domain.TaskSearchService
	pb.UnimplementedTasksServer
}

func New(domain domain.TaskService, search domain.TaskSearchService) *Server {

	s := &Server{domain: domain, search: search}

	// grpc server
	gs, err := kitGrpc.NewGrpcServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterTasksServer(s.Srv, s)

	return s
}

func (s *Server) ListenAsync() {

	go func () {
		err := s.Server.Listen("localhost", "50052")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) New(ctx context.Context, rq *pb.NewTaskRequest) (*pb.Task, error) {

	task, err := s.domain.New(s.fromPb(rq))
	if err != nil {
		return nil, err
	}

	return s.fromDomain(task), nil
}

func (s *Server) NextTransitions(ctx context.Context, rq *pb.NextTransitionsRequest) (*pb.NextTransitionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NextTransitions not implemented")
}

func (s *Server) MakeTransition(ctx context.Context, rq *pb.MakeTransitionRequest) (*pb.Task, error) {

	task, err := s.domain.MakeTransition(rq.TaskId, rq.TransitionId)
	if err != nil {
		return nil, err
	}

	return s.fromDomain(task), nil
}

func (s *Server) GetByChannel(ctx context.Context, rq *pb.GetByChannelRequest) (*pb.GetByChannelResponse, error) {

	response := &pb.GetByChannelResponse{Tasks: []*pb.Task{}}

	tasks := s.domain.GetByChannel(rq.ChannelId)
	for _, t := range tasks {
		response.Tasks = append(response.Tasks, s.fromDomain(t))
	}

	return response, nil

}

func (s *Server) SetAssignee(ctx context.Context, rq *pb.SetAssigneeRequest) (*pb.Task, error) {

	task, err := s.domain.SetAssignee(rq.TaskId, s.assigneeFromPb(rq.Assignee))
	if err != nil {
		return nil, err
	}

	return s.fromDomain(task), nil
}

func (s *Server) GetById(ctx context.Context, rq *pb.GetByIdRequest) (*pb.Task, error) {
	task := s.domain.Get(rq.Id)
	return s.fromDomain(task), nil
}

func (s *Server) Search(ctx context.Context, rq *pb.SearchRequest) (*pb.SearchResponse, error) {

	dRs, err := s.search.Search(s.searchRqFromPb(rq))
	if err != nil {
		return nil, err
	}

	return s.searchRsFromDomain(dRs), nil
}

func (s *Server) GetAssignmentLog(ctx context.Context, rq *pb.AssignmentLogRequest) (*pb.AssignmentLogResponse, error) {
	dRs, err := s.domain.GetAssignmentLog(s.assLogRqFromPb(rq))
	if err != nil {
		return nil, err
	}

	return s.assLogRsFromDomain(dRs), nil
}

func (s *Server) GetHistory(ctx context.Context, rq *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {

	rs := &pb.GetHistoryResponse{Items: []*pb.History{}}

	for _, h := range s.domain.GetHistory(rq.TaskId) {
		rs.Items = append(rs.Items, s.histToPb(h))
	}

	return rs, nil
}