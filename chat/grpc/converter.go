package grpc

import (
	"gitlab.medzdrav.ru/prototype/chat/domain"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

func (s *Server) toFromDomain(r *pb.From) *domain.From {
	return &domain.From{
		Who:        domain.Who(r.Who),
		ChatUserId: r.ChatUserId,
	}
}
