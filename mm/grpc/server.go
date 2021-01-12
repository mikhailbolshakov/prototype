package grpc

import (
	"context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/mm/domain"
	pb "gitlab.medzdrav.ru/prototype/proto/mm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	*kitGrpc.Server
	domain domain.MMService
	pb.UnimplementedUsersServer
	pb.UnimplementedChannelsServer
	pb.UnimplementedPostsServer
}

func New(domain domain.MMService) *Server {

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

	go func () {
		err := s.Server.Listen("localhost", "50053")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) CreateUser(ctx context.Context, rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}

func (s *Server) CreateClientChannel(ctx context.Context, rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateClientChannel not implemented")
}
