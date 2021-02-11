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

func (s *serviceImpl) StartProcess(ctx context.Context, rq *pb.StartProcessRequest) (*pb.StartProcessResponse, error) {
	return s.ProcessClient.StartProcess(ctx, rq)
}

