package grpc

import (
	"context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/mm/domain"
	pb "gitlab.medzdrav.ru/prototype/proto/mm"
	"log"
)

type Server struct {
	*kitGrpc.Server
	domain domain.Service
	pb.UnimplementedUsersServer
	pb.UnimplementedChannelsServer
	pb.UnimplementedPostsServer
}

func New(domain domain.Service) *Server {

	s := &Server{domain: domain}

	// grpc server
	gs, err := kitGrpc.NewGrpcServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterUsersServer(s.Srv, s)
	pb.RegisterChannelsServer(s.Srv, s)
	pb.RegisterPostsServer(s.Srv, s)

	return s
}

func (s *Server) ListenAsync() {

	go func() {
		err := s.Server.Listen("localhost", "50053")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) CreateUser(ctx context.Context, rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	rs, err := s.domain.CreateUser(&domain.CreateUserRequest{
		Username: rq.Username,
		Email:    rq.Email,
	})
	if err != nil {
		return nil, err
	}
	response := &pb.CreateUserResponse{Id: rs.Id}

	return response, nil
}

func (s *Server) CreateClientChannel(ctx context.Context, rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error) {

	rs, err := s.domain.CreateClientChannel(&domain.CreateClientChannelRequest{
		ClientUserId: rq.ClientUserId,
		DisplayName:  rq.DisplayName,
		Name:         rq.Name,
		Subscribers:  rq.Subscribers,
	})
	if err != nil {
		return nil, err
	}
	response := &pb.CreateClientChannelResponse{ChannelId: rs.ChannelId}

	return response, nil
}

func (s *Server) GetUsersStatuses(ctx context.Context, rq *pb.GetUsersStatusesRequest) (*pb.GetUserStatusesResponse, error) {

	rs, err := s.domain.GetUsersStatuses(&domain.GetUsersStatusesRequest{UserIds: rq.MMUserIds})
	if err != nil {
		return nil, err
	}
	response := &pb.GetUserStatusesResponse{Statuses: []*pb.UserStatus{}}
	for _, s := range rs.Statuses {
		response.Statuses = append(response.Statuses, &pb.UserStatus{
			Status:   s.Status,
			MMUserId: s.UserId,
		})
	}
	return response, nil

}
