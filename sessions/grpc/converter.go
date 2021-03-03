package grpc

import (
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/sessions"
	"gitlab.medzdrav.ru/prototype/sessions/domain"
)

func (s *Server) toSessionPb(r *domain.Session) *pb.Session {
	if r == nil {
		return nil
	}
	return &pb.Session{
		Id:            r.Id,
		UserId:        r.UserId,
		Username:      r.Username,
		ChatUserId:    r.ChatUserId,
		ChatSessionId: r.ChatSessionId,
		LoginAt:       grpc.TimeToPbTS(&r.LoginAt),
	}
}

func (s *Server) toSessionsPb(r []*domain.Session) []*pb.Session {
	var rs  []*pb.Session
	for _, ss := range r{
		rs = append(rs, s.toSessionPb(ss))
	}
	return rs
}

func (s *Server) toUserSessionsPb(us *domain.UserSessionInfo) *pb.UserSessionsInfo {
	res := &pb.UserSessionsInfo{
		UserId:     us.UserId,
		ChatUserId: us.ChatUserId,
		Sessions:   []*pb.SessionInfo{},
	}

	for _, s := range us.Sessions {
		res.Sessions = append(res.Sessions, &pb.SessionInfo{
			Id:             s.Id,
			StartAt:        grpc.TimeToPbTS(&s.StartAt),
			SentWsMessages: s.SentWsMessages,
			ChatSessionId:  s.ChatSessionId,
		})
	}

	return res
}