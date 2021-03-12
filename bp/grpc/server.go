package grpc

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/bp/logger"
	"gitlab.medzdrav.ru/prototype/bp/meta"
	bpmKit "gitlab.medzdrav.ru/prototype/kit/bpm"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/bp"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"log"
)

type Server struct {
	host, port string
	*kitGrpc.Server
	bpm bpmKit.Engine
	pb.UnimplementedProcessServer
}

func New(bpm bpmKit.Engine) *Server {

	s := &Server{bpm: bpm}

	// grpc server
	gs, err := kitGrpc.NewServer(meta.Meta.ServiceCode(), logger.LF())
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterProcessServer(s.Srv, s)

	return s
}

func  (s *Server) Init(c *config.Config) error {
	cfg := c.Services[meta.Meta.ServiceCode()]
	s.host = cfg.Grpc.Host
	s.port = cfg.Grpc.Port
	return nil
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
		err := s.Server.Listen(s.host, s.port)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

