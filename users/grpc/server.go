package grpc

import (
	"context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"gitlab.medzdrav.ru/prototype/users/logger"
	"gitlab.medzdrav.ru/prototype/users/meta"
	"log"
)

type Server struct {
	port, host string
	*kitGrpc.Server
	domain domain.UserService
	pb.UnimplementedUsersServer
}

func New(domain domain.UserService) *Server {

	s := &Server{domain: domain}

	// grpc server
	gs, err := kitGrpc.NewServer(meta.Meta.ServiceCode(), logger.LF())
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterUsersServer(s.Srv, s)

	return s

}

func  (s *Server) Init(c *config.Config) error {
	usersCfg := c.Services["users"]
	s.host = usersCfg.Grpc.Host
	s.port = usersCfg.Grpc.Port
	return nil
}

func (s *Server) ListenAsync() {

	go func() {
		err := s.Server.Listen(s.host, s.port)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) CreateClient(ctx context.Context, rq *pb.CreateClientRequest) (*pb.User, error) {

	client := &domain.User{
		Type: domain.USER_TYPE_CLIENT,
		ClientDetails: &domain.ClientDetails{
			FirstName:         rq.FirstName,
			MiddleName:        rq.MiddleName,
			LastName:          rq.LastName,
			Sex:               rq.Sex,
			BirthDate:         *(kitGrpc.PbTSToTime(rq.BirthDate)),
			Phone:             rq.Phone,
			Email:             rq.Email,
			PersonalAgreement: &domain.PersonalAgreement{},
			PhotoUrl:          rq.PhotoUrl,
		},
	}
	user, err := s.domain.Create(ctx, client)
	if err != nil {
		return nil, err
	}
	return s.toUserPb(user), nil
}

func (s *Server) CreateConsultant(ctx context.Context, rq *pb.CreateConsultantRequest) (*pb.User, error) {
	consultant := &domain.User{
		Type: domain.USER_TYPE_CONSULTANT,
		ConsultantDetails: &domain.ConsultantDetails{
			FirstName:  rq.FirstName,
			MiddleName: rq.MiddleName,
			LastName:   rq.LastName,
			Email:      rq.Email,
			PhotoUrl:   rq.PhotoUrl,
		},
		Groups: rq.Groups,
	}

	user, err := s.domain.Create(ctx, consultant)
	if err != nil {
		return nil, err
	}
	return s.toUserPb(user), nil
}

func (s *Server) CreateExpert(ctx context.Context, rq *pb.CreateExpertRequest) (*pb.User, error) {
	expert := &domain.User{
		Type: domain.USER_TYPE_EXPERT,
		ExpertDetails: &domain.ExpertDetails{
			FirstName:  rq.FirstName,
			MiddleName: rq.MiddleName,
			LastName:   rq.LastName,
			Email:      rq.Email,
			PhotoUrl:   rq.PhotoUrl,
		},
		Groups: rq.Groups,
	}

	user, err := s.domain.Create(ctx, expert)
	if err != nil {
		return nil, err
	}
	return s.toUserPb(user), nil
}

func (s *Server) GetByUsername(ctx context.Context, rq *pb.GetByUsernameRequest) (*pb.User, error) {
	user := s.domain.GetByUsername(ctx, rq.Username)
	return s.toUserPb(user), nil
}

func (s *Server) GetByMMId(ctx context.Context, rq *pb.GetByMMIdRequest) (*pb.User, error) {
	user := s.domain.GetByMMId(ctx, rq.MMId)
	return s.toUserPb(user), nil
}

func (s *Server) Get(ctx context.Context, rq *pb.GetByIdRequest) (*pb.User, error) {
	user := s.domain.Get(ctx, rq.Id)
	return s.toUserPb(user), nil
}

func (s *Server) Search(ctx context.Context, rq *pb.SearchRequest) (*pb.SearchResponse, error) {

	dRs, err := s.domain.Search(ctx, s.toSrchRqDomain(rq))
	if err != nil {
		return nil, err
	}

	return s.toSrchRsPb(dRs), nil
}

func (s *Server) Activate(ctx context.Context, rq *pb.ActivateRequest) (*pb.User, error) {
	user, err := s.domain.Activate(ctx, rq.UserId)
	if err != nil {
		return nil, err
	}
	return s.toUserPb(user), nil
}

func (s *Server) Delete(ctx context.Context, rq *pb.DeleteRequest) (*pb.User, error) {
	user, err := s.domain.Delete(ctx, rq.UserId)
	if err != nil {
		return nil, err
	}
	return s.toUserPb(user), nil
}

func (s *Server) SetClientDetails(ctx context.Context, rq *pb.SetClientDetailsRequest) (*pb.User, error) {

	user, err := s.domain.SetClientDetails(ctx, rq.UserId, s.toClientDetailsDomain(rq.ClientDetails))
	if err != nil {
		return nil, err
	}
	return s.toUserPb(user), nil
}

func (s *Server) SetMMUserId(ctx context.Context, rq *pb.SetMMIdRequest) (*pb.User, error) {
	user, err := s.domain.SetMMUserId(ctx, rq.UserId, rq.MMId)
	if err != nil {
		return nil, err
	}
	return s.toUserPb(user), nil
}

func (s *Server) SetKKUserId(ctx context.Context, rq *pb.SetKKIdRequest) (*pb.User, error) {
	user, err := s.domain.SetKKUserId(ctx, rq.UserId, rq.KKId)
	if err != nil {
		return nil, err
	}
	return s.toUserPb(user), nil
}
