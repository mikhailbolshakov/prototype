package grpc

import (
	"context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/sessions"
	"gitlab.medzdrav.ru/prototype/sessions/domain"
	"gitlab.medzdrav.ru/prototype/sessions/logger"
	"gitlab.medzdrav.ru/prototype/sessions/meta"
)

type Server struct {
	host, port string
	*kitGrpc.Server
	domain  domain.SessionsService
	monitor domain.SessionMonitor
	pb.UnimplementedSessionsServer
	pb.UnimplementedMonitorServer
}

func New(domain domain.SessionsService, monitor domain.SessionMonitor) *Server {

	s := &Server{domain: domain, monitor: monitor}

	// grpc server
	gs, err := kitGrpc.NewServer(meta.Meta.ServiceCode(), logger.LF())
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterSessionsServer(s.Srv, s)
	pb.RegisterMonitorServer(s.Srv, s)

	return s
}

func (s *Server) Init(c *config.Config) error {
	cfg := c.Services["sessions"]
	s.host = cfg.Grpc.Host
	s.port = cfg.Grpc.Port
	return nil
}

func (s *Server) ListenAsync() {

	go func() {
		err := s.Server.Listen(s.host, s.port)
		if err != nil {
			return
		}
	}()
}

func (s *Server) Login(ctx context.Context, rq *pb.LoginRequest) (*pb.LoginResponse, error) {

	rs, err := s.domain.Login(ctx, &domain.LoginRequest{
		Username:  rq.Username,
		Password:  rq.Password,
		ChatLogin: rq.ChatLogin,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{SessionId: rs.SessionId}, nil

}
func (s *Server) Logout(ctx context.Context, rq *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := s.domain.Logout(ctx, &domain.LogoutRequest{UserId: rq.UserId})
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResponse{}, nil
}
func (s *Server) Get(ctx context.Context, rq *pb.GetByIdRequest) (*pb.Session, error) {
	ss, err := s.domain.Get(ctx, rq.Id)
	if err != nil {
		return nil, err
	}
	return s.toSessionPb(ss), nil
}
func (s *Server) GetByUser(ctx context.Context, rq *pb.GetByUserRequest) (*pb.SessionsResponse, error) {
	rs, err := s.domain.GetByUser(ctx, &domain.GetByUserRequest{
		UserId:   rq.UserId,
		Username: rq.Username,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SessionsResponse{Sessions: s.toSessionsPb(rs)}, nil
}

func (s *Server) AuthSession(ctx context.Context, rq *pb.AuthSessionRequest) (*pb.Session, error) {
	rs, err := s.domain.AuthSession(ctx, rq.SessionId)
	if err != nil {
		return nil, err
	}
	return s.toSessionPb(rs), nil
}

func (s *Server) UserSessions(ctx context.Context, rq *pb.UserSessionsRequest) (*pb.UserSessionsInfo, error) {
	ss := s.monitor.GetUserSessions(ctx, rq.UserId)
	return s.toUserSessionsPb(ss), nil
}

func (s *Server) TotalSessions(ctx context.Context, rq *pb.SessionsTotalRequest) (*pb.TotalSessionInfo, error) {
	si := s.monitor.GetTotalSessions(ctx)
	return &pb.TotalSessionInfo{
		ActiveCount:      uint32(si.ActiveCount),
		ActiveUsersCount: uint32(si.ActiveUsersCount),
	}, nil
}
