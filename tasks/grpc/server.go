package grpc

import (
	"context"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	log "gitlab.medzdrav.ru/prototype/kit/log"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	domain "gitlab.medzdrav.ru/prototype/tasks/domain"
	"gitlab.medzdrav.ru/prototype/tasks/meta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	host, port string
	*kitGrpc.Server
	domain domain.TaskService
	pb.UnimplementedTasksServer
}

func New(domain domain.TaskService) *Server {

	s := &Server{domain: domain}

	// grpc server
	gs, err := kitGrpc.NewServer(meta.ServiceCode)
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterTasksServer(s.Srv, s)

	return s
}

func  (s *Server) Init(c *kitConfig.Config) error {
	usersCfg := c.Services["tasks"]
	s.host = usersCfg.Grpc.Host
	s.port = usersCfg.Grpc.Port
	return nil
}

func (s *Server) ListenAsync() {

	go func () {
		err := s.Server.Listen(s.host, s.port)
		if err != nil {
			log.L().Pr("grpc").Mth("listen").E(err).Err()
			return
		}
	}()
}

func (s *Server) New(ctx context.Context, rq *pb.NewTaskRequest) (*pb.Task, error) {

	task, err := s.domain.New(ctx, s.toTaskDomain(rq))
	if err != nil {
		return nil, err
	}

	return s.toTaskPb(task), nil
}

func (s *Server) NextTransitions(ctx context.Context, rq *pb.NextTransitionsRequest) (*pb.NextTransitionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NextTransitions not implemented")
}

func (s *Server) MakeTransition(ctx context.Context, rq *pb.MakeTransitionRequest) (*pb.Task, error) {

	task, err := s.domain.MakeTransition(ctx, rq.TaskId, rq.TransitionId)
	if err != nil {
		return nil, err
	}

	return s.toTaskPb(task), nil
}

func (s *Server) GetByChannel(ctx context.Context, rq *pb.GetByChannelRequest) (*pb.GetByChannelResponse, error) {

	response := &pb.GetByChannelResponse{Tasks: []*pb.Task{}}

	tasks := s.domain.GetByChannel(ctx, rq.ChannelId)
	for _, t := range tasks {
		response.Tasks = append(response.Tasks, s.toTaskPb(t))
	}

	return response, nil

}

func (s *Server) SetAssignee(ctx context.Context, rq *pb.SetAssigneeRequest) (*pb.Task, error) {

	task, err := s.domain.SetAssignee(ctx, rq.TaskId, s.toAssigneeDomain(rq.Assignee))
	if err != nil {
		return nil, err
	}

	return s.toTaskPb(task), nil
}

func (s *Server) GetById(ctx context.Context, rq *pb.GetByIdRequest) (*pb.Task, error) {
	task := s.domain.Get(ctx, rq.Id)
	return s.toTaskPb(task), nil
}

func (s *Server) Search(ctx context.Context, rq *pb.SearchRequest) (*pb.SearchResponse, error) {

	dRs, err := s.domain.Search(ctx, s.toSrchRqDomain(rq))
	if err != nil {
		return nil, err
	}

	return s.toSrchRsPb(dRs), nil
}

func (s *Server) GetAssignmentLog(ctx context.Context, rq *pb.AssignmentLogRequest) (*pb.AssignmentLogResponse, error) {
	dRs, err := s.domain.GetAssignmentLog(ctx, s.toAssignLogDomain(rq))
	if err != nil {
		return nil, err
	}

	return s.toAssignLogRsPb(dRs), nil
}

func (s *Server) GetHistory(ctx context.Context, rq *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {

	rs := &pb.GetHistoryResponse{Items: []*pb.History{}}

	for _, h := range s.domain.GetHistory(ctx, rq.TaskId) {
		rs.Items = append(rs.Items, s.toHistoryPb(h))
	}

	return rs, nil
}

