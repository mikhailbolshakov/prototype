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

func (s *serviceImpl) Login(ctx context.Context, rq *pb.LoginRequest) (string, error) {
	rs, err := s.SessionsClient.Login(ctx, rq)
	if err != nil {
		return "", err
	}
	return rs.SessionId, nil
}

func (s *serviceImpl) Logout(ctx context.Context, userId string) error {
	_, err := s.SessionsClient.Logout(ctx, &pb.LogoutRequest{UserId: userId})
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceImpl) Get(ctx context.Context, sid string) (*pb.Session, error) {
	ss, err := s.SessionsClient.Get(ctx, &pb.GetByIdRequest{Id: sid})
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (s *serviceImpl) GetByUser(ctx context.Context, userId, username string) ([]*pb.Session, error) {
	rs, err := s.SessionsClient.GetByUser(ctx, &pb.GetByUserRequest{
		UserId:   userId,
		Username: username,
	})
	if err != nil {
		return nil, err
	}
	return rs.Sessions, nil
}

func (s *serviceImpl) AuthSession(ctx context.Context, sid string) (*pb.Session, error) {
	rs, err := s.SessionsClient.AuthSession(ctx, &pb.AuthSessionRequest{SessionId: sid})
	if err != nil {
		return nil, err
	}
	return rs, nil
}