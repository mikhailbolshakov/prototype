package grpc

import (
	"context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"log"
)

type Server struct {
	*kitGrpc.Server
	domain domain.UserService
	search domain.UserSearchService
	pb.UnimplementedUsersServer
}

func New(domain domain.UserService, search domain.UserSearchService) *Server {

	s := &Server{domain: domain, search: search}

	// grpc server
	gs, err := kitGrpc.NewGrpcServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterUsersServer(s.Srv, s)

	return s

}

func (s *Server) ListenAsync() {

	go func () {
		err := s.Server.Listen("localhost", "50051")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Create(ctx context.Context, rq *pb.CreateUserRequest) (*pb.User, error) {

	user, err := s.domain.Create(s.fromPb(rq))
	if err != nil {
		return nil, err
	}

	return s.fromDomain(user), nil
}

func (s *Server) GetByUsername(ctx context.Context, rq *pb.GetByUsernameRequest) (*pb.User, error) {

	user := s.domain.GetByUsername(rq.Username)
	return s.fromDomain(user), nil
}

func (s *Server) GetByMMId(ctx context.Context, rq *pb.GetByMMIdRequest) (*pb.User, error) {

	user := s.domain.GetByMMId(rq.MMId)
	return s.fromDomain(user), nil
}

func (s *Server) Search(ctx context.Context, rq *pb.SearchRequest) (*pb.SearchResponse, error) {

	dRs, err := s.search.Search(s.searchRqFromPb(rq))
	if err != nil {
		return nil, err
	}

	return s.searchRsFromDomain(dRs), nil
}