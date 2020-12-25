package users

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/users/grpc"
	"gitlab.medzdrav.ru/prototype/users/repository"
	"log"
)

type Service struct {
	*kitGrpc.Server
}

func NewService() *Service {

	s := &Service{}

	// grpc server
	gs, err := kitGrpc.NewGrpcServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterUsersServer(s.Srv, grpc.New())

	return s
}

func (s *Service) Start() error {

	// init infrastructure
	err := repository.InitInfrastructure()
	if err != nil {
		return err
	}

	// grpc server
	go func () {
		err := s.Listen("localhost", "50051")
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}