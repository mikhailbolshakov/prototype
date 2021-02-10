package grpc

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/users/domain"
)

func (s *Server) toClientDetailsDomain(c *pb.ClientDetails) *domain.ClientDetails {
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
		CommonChannelId: c.CommonChannelId,
		MedChannelId:    c.MedChannelId,
		LawChannelId:    c.LawChannelId,
		PhotoUrl: c.PhotoUrl,
	}
}

func (s *Server) toUserPb(user *domain.User) *pb.User {

	if user == nil {
		return nil
	}

	pu := &pb.User{
		Id:       user.Id,
		Type:     user.Type,
		Username: user.Username,
		Status:   user.Status,
		MMId:     user.MMUserId,
		Groups:   user.Groups,
	}

	switch user.Type {
	case domain.USER_TYPE_CLIENT:
		pu.ClientDetails = &pb.ClientDetails{
			FirstName:       user.ClientDetails.FirstName,
			MiddleName:      user.ClientDetails.MiddleName,
			LastName:        user.ClientDetails.LastName,
			Sex:             user.ClientDetails.Sex,
			BirthDate:       grpc.TimeToPbTS(&user.ClientDetails.BirthDate),
			Phone:           user.ClientDetails.Phone,
			Email:           user.ClientDetails.Email,
			CommonChannelId: user.ClientDetails.CommonChannelId,
			MedChannelId:    user.ClientDetails.MedChannelId,
			LawChannelId:    user.ClientDetails.LawChannelId,
			PersonalAgreement: &pb.PersonalAgreement{
				GivenAt:   grpc.TimeToPbTS(user.ClientDetails.PersonalAgreement.GivenAt),
				RevokedAt: grpc.TimeToPbTS(user.ClientDetails.PersonalAgreement.RevokedAt),
			},
			PhotoUrl: user.ClientDetails.PhotoUrl,
		}
	case domain.USER_TYPE_CONSULTANT:
		pu.ConsultantDetails = &pb.ConsultantDetails{
			FirstName:  user.ConsultantDetails.FirstName,
			MiddleName: user.ConsultantDetails.MiddleName,
			LastName:   user.ConsultantDetails.LastName,
			Email:      user.ConsultantDetails.Email,
			PhotoUrl: user.ConsultantDetails.PhotoUrl,
		}
	case domain.USER_TYPE_EXPERT:
		pu.ExpertDetails = &pb.ExpertDetails{
			FirstName:      user.ExpertDetails.FirstName,
			MiddleName:     user.ExpertDetails.MiddleName,
			LastName:       user.ExpertDetails.LastName,
			Email:          user.ExpertDetails.Email,
			PhotoUrl: user.ExpertDetails.PhotoUrl,
		}
	case domain.USER_TYPE_SUPERVISOR:
	}

	return pu

}

func (s *Server) toSrchRqDomain(pb *pb.SearchRequest) *domain.SearchCriteria {

	if pb == nil {
		return nil
	}

	return &domain.SearchCriteria{
		PagingRequest: &common.PagingRequest{
			Size:  int(pb.Paging.Size),
			Index: int(pb.Paging.Index),
		},
		UserType:        pb.UserType,
		Username:        pb.Username,
		UserGroup:       pb.UserGroup,
		Status:          pb.Status,
		Email:           pb.Email,
		Phone:           pb.Phone,
		MMId:            pb.MMId,
		CommonChannelId: pb.CommonChannelId,
		MedChannelId:    pb.MedChannelId,
		LawChannelId:    pb.LawChannelId,
		OnlineStatuses:  pb.OnlineStatuses,
	}

}

func (s *Server) toSrchRsPb(d *domain.SearchResponse) *pb.SearchResponse {

	rs := &pb.SearchResponse{
		Paging: &pb.PagingResponse{
			Total: int32(d.PagingResponse.Total),
			Index: int32(d.PagingResponse.Index),
		},
		Users: []*pb.User{},
	}

	for _, t := range d.Users {
		rs.Users = append(rs.Users, s.toUserPb(t))
	}

	return rs
}
