package grpc

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/users/domain"
)

func (s *Server) clientDetailsFromPb(c *pb.ClientDetails) *domain.ClientDetails {
	return &domain.ClientDetails{
		FirstName:  c.FirstName,
		MiddleName: c.MiddleName,
		LastName:   c.LastName,
		Sex:        c.Sex,
		BirthDate:  *(grpc.PbTSToTime(c.BirthDate)),
		Phone:      c.Phone,
		Email:      c.Email,
		PersonalAgreement: &domain.PersonalAgreement{
			GivenAt:   grpc.PbTSToTime(c.PersonalAgreement.GivenAt),
			RevokedAt: grpc.PbTSToTime(c.PersonalAgreement.RevokedAt),
		},
		MMChannelId: c.MMChannelId,
	}
}

func (s *Server) fromDomain(user *domain.User) *pb.User {

	if user == nil {
		return nil
	}

	pu := &pb.User{
		Id:       user.Id,
		Type:     user.Type,
		Username: user.Username,
		Status:   user.Status,
		MMId:     user.MMUserId,
	}

	switch user.Type {
	case domain.USER_TYPE_CLIENT:
		pu.ClientDetails = &pb.ClientDetails{
			FirstName:   user.ClientDetails.FirstName,
			MiddleName:  user.ClientDetails.MiddleName,
			LastName:    user.ClientDetails.LastName,
			Sex:         user.ClientDetails.Sex,
			BirthDate:   grpc.TimeToPbTS(&user.ClientDetails.BirthDate),
			Phone:       user.ClientDetails.Phone,
			Email:       user.ClientDetails.Email,
			MMChannelId: user.ClientDetails.MMChannelId,
			PersonalAgreement: &pb.PersonalAgreement{
				GivenAt:   grpc.TimeToPbTS(user.ClientDetails.PersonalAgreement.GivenAt),
				RevokedAt: grpc.TimeToPbTS(user.ClientDetails.PersonalAgreement.RevokedAt),
			},
		}
	case domain.USER_TYPE_CONSULTANT:
		pu.ConsultantDetails = &pb.ConsultantDetails{
			FirstName:  user.ConsultantDetails.FirstName,
			MiddleName: user.ConsultantDetails.MiddleName,
			LastName:   user.ConsultantDetails.LastName,
			Email:      user.ConsultantDetails.Email,
		}
	case domain.USER_TYPE_EXPERT:
		pu.ExpertDetails = &pb.ExpertDetails{
			FirstName:      user.ExpertDetails.FirstName,
			MiddleName:     user.ExpertDetails.MiddleName,
			LastName:       user.ExpertDetails.LastName,
			Email:          user.ExpertDetails.Email,
			Specialization: user.ExpertDetails.Specialization,
		}
	case domain.USER_TYPE_SUPERVISOR:
	}

	return pu

}

func (s *Server) searchRqFromPb(pb *pb.SearchRequest) *domain.SearchCriteria {

	if pb == nil {
		return nil
	}

	return &domain.SearchCriteria{
		PagingRequest: &common.PagingRequest{
			Size:  int(pb.Paging.Size),
			Index: int(pb.Paging.Index),
		},
		UserType:       pb.UserType,
		Username:       pb.Username,
		Status:         pb.Status,
		Email:          pb.Email,
		Phone:          pb.Phone,
		MMId:           pb.MMId,
		MMChannelId:    pb.MMChannelId,
		OnlineStatuses: pb.OnlineStatuses,
	}

}

func (s *Server) searchRsFromDomain(d *domain.SearchResponse) *pb.SearchResponse {

	rs := &pb.SearchResponse{
		Paging: &pb.PagingResponse{
			Total: int32(d.PagingResponse.Total),
			Index: int32(d.PagingResponse.Index),
		},
		Users: []*pb.User{},
	}

	for _, t := range d.Users {
		rs.Users = append(rs.Users, s.fromDomain(t))
	}

	return rs
}
