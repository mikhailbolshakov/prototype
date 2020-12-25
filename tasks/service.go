package tasks

import (
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/tasks/grpc"
	"gitlab.medzdrav.ru/prototype/tasks/repository"
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
	pb.RegisterTasksServer(s.Srv, grpc.New())

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
		err := s.Listen("localhost", "50052")
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}