package grpc

import (
	"context"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"log"
)

type Server struct {
	port, host string
	*kitGrpc.Server
	domain domain.UserService
	search domain.UserSearchService
	pb.UnimplementedUsersServer
}

func New(domain domain.UserService, search domain.UserSearchService) *Server {

	s := &Server{domain: domain, search: search}

	// grpc server
	gs, err := kitGrpc.NewServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterUsersServer(s.Srv, s)

	return s

}

func  (s *Server) Init(c *kitConfig.Config) error {
	usersCfg := c.Services["users"]
	s.host = usersCfg.Grpc.Hosts[0]
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
	user, err := s.domain.Create(client)
	if err != nil {
		return nil, err
	}
	return s.fromDomain(user), nil
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

	user, err := s.domain.Create(consultant)
	if err != nil {
		return nil, err
	}
	return s.fromDomain(user), nil
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

	user, err := s.domain.Create(expert)
	if err != nil {
		return nil, err
	}
	return s.fromDomain(user), nil
}

func (s *Server) GetByUsername(ctx context.Context, rq *pb.GetByUsernameRequest) (*pb.User, error) {
	user := s.domain.GetByUsername(rq.Username)
	return s.fromDomain(user), nil
}

func (s *Server) GetByMMId(ctx context.Context, rq *pb.GetByMMIdRequest) (*pb.User, error) {
	user := s.domain.GetByMMId(rq.MMId)
	return s.fromDomain(user), nil
}

func (s *Server) Get(ctx context.Context, rq *pb.GetByIdRequest) (*pb.User, error) {
	user := s.domain.Get(rq.Id)
	return s.fromDomain(user), nil
}

func (s *Server) Search(ctx context.Context, rq *pb.SearchRequest) (*pb.SearchResponse, error) {

	dRs, err := s.search.Search(s.searchRqFromPb(rq))
	if err != nil {
		return nil, err
	}

	return s.searchRsFromDomain(dRs), nil
}

func (s *Server) Activate(ctx context.Context, rq *pb.ActivateRequest) (*pb.User, error) {
	user, err := s.domain.Activate(rq.UserId)
	if err != nil {
		return nil, err
	}
	return s.fromDomain(user), nil
}

func (s *Server) Delete(ctx context.Context, rq *pb.DeleteRequest) (*pb.User, error) {
	user, err := s.domain.Delete(rq.UserId)
	if err != nil {
		return nil, err
	}
	return s.fromDomain(user), nil
}

func (s *Server) SetClientDetails(ctx context.Context, rq *pb.SetClientDetailsRequest) (*pb.User, error) {

	user, err := s.domain.SetClientDetails(rq.UserId, s.clientDetailsFromPb(rq.ClientDetails))
	if err != nil {
		return nil, err
	}
	return s.fromDomain(user), nil
}

func (s *Server) SetMMUserId(ctx context.Context, rq *pb.SetMMIdRequest) (*pb.User, error) {
	user, err := s.domain.SetMMUserId(rq.UserId, rq.MMId)
	if err != nil {
		return nil, err
	}
	return s.fromDomain(user), nil
}

func (s *Server) SetKKUserId(ctx context.Context, rq *pb.SetKKIdRequest) (*pb.User, error) {
	user, err := s.domain.SetKKUserId(rq.UserId, rq.KKId)
	if err != nil {
		return nil, err
	}
	return s.fromDomain(user), nil
}
