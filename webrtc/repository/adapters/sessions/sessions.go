package sessions

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/sessions"
)

type serviceImpl struct {
	pb.SessionsClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{
	}
	return a
}

func (s *serviceImpl) AuthSession(ctx context.Context, sid string) (*pb.Session, error) {
	rs, err := s.SessionsClient.AuthSession(ctx, &pb.AuthSessionRequest{SessionId: sid})
	if err != nil {
		return nil, err
	}
	return rs, nil
}