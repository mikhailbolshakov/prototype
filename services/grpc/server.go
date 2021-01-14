package grpc

import (
	"context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	*kitGrpc.Server
	balanceService domain.UserBalanceService
	pb.UnimplementedUserServicesServer
}

func New(domain domain.UserBalanceService) *Server {

	s := &Server{balanceService: domain}

	// grpc server
	gs, err := kitGrpc.NewGrpcServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterUserServicesServer(s.Srv, s)

	return s
}

func (s *Server) ListenAsync() {

	go func () {
		err := s.Server.Listen("localhost", "50054")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Add(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {

	rs, err := s.balanceService.Add(&domain.ModifyBalanceRequest{
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      int(rq.Quantity),
	})
	if err != nil {
		return nil, err
	}

	return userBalanceFromDomain(rs), nil

}

func (s *Server) GetBalance(ctx context.Context, rq *pb.GetBalanceRequest) (*pb.UserBalance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (s *Server) WriteOff(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteOff not implemented")
}
func (s *Server) Lock(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lock not implemented")
}