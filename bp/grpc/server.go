package grpc

import (
	"context"
	"encoding/json"
	bpmKit "gitlab.medzdrav.ru/prototype/kit/bpm"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/bp"
	"log"
)

type Server struct {
	*kitGrpc.Server
	bpm bpmKit.Engine
	pb.UnimplementedProcessServer
}

func New(bpm bpmKit.Engine) *Server {

	s := &Server{bpm: bpm}

	// grpc server
	gs, err := kitGrpc.NewServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterProcessServer(s.Srv, s)

	return s
}

func (s *Server) StartProcess(ctx context.Context, rq *pb.StartProcessRequest) (*pb.StartProcessResponse, error) {

	vars := map[string]interface{}{}
	if rq.Vars != nil {
		_ = json.Unmarshal(rq.Vars, &vars)
	}

	rs, err := s.bpm.StartProcess(rq.ProcessId, vars)
	if err != nil {
		return nil, err
	}
	return &pb.StartProcessResponse{Id: rs}, nil

}

func (s *Server) ListenAsync() {

	go func() {
		err := s.Server.Listen("localhost", "50055")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

