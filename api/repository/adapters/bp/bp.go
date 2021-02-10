package bp

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/bp"
)

type serviceImpl struct {
	pb.ProcessClient
}

func newServiceImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (s *serviceImpl) StartProcess(rq *pb.StartProcessRequest) (*pb.StartProcessResponse, error) {
	return s.ProcessClient.StartProcess(context.Background(), rq)
}

