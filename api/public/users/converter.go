package users

import (
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

const (
	USER_TYPE_CLIENT     = "client"
	USER_TYPE_CONSULTANT = "consultant"
	USER_TYPE_EXPERT     = "expert"
	USER_TYPE_SUPERVISOR = "supervisor"
)

func (s *ctrlImpl) fromPb(r *pb.User) *User {
	user := &User{
		Id:       r.Id,
		Username: r.Username,
		Type:     r.Type,
		Status:   r.Status,
		MMUserId: r.MMId,
		KKUserId: r.KKId,
	}

	switch user.Type {
	case USER_TYPE_CLIENT:
		user.ClientDetails = &ClientDetails{
			FirstName:  r.ClientDetails.FirstName,
			MiddleName: r.ClientDetails.MiddleName,
			LastName:   r.ClientDetails.LastName,
			Sex:        r.ClientDetails.Sex,
			BirthDate:  *(grpc.PbTSToTime(r.ClientDetails.BirthDate)),
			Phone:      r.ClientDetails.Phone,
			Email:      r.ClientDetails.Email,
			PersonalAgreement: &PersonalAgreement{
				GivenAt:   grpc.PbTSToTime(r.ClientDetails.PersonalAgreement.GivenAt),
				RevokedAt: grpc.PbTSToTime(r.ClientDetails.PersonalAgreement.RevokedAt),
			},
			MMChannelId: r.ClientDetails.MMChannelId,
		}
	case USER_TYPE_CONSULTANT:
		user.ConsultantDetails = &ConsultantDetails{
			FirstName:  r.ConsultantDetails.FirstName,
			MiddleName: r.ConsultantDetails.MiddleName,
			LastName:   r.ConsultantDetails.LastName,
			Email:      r.ConsultantDetails.Email,
		}
	case USER_TYPE_EXPERT:
		user.ExpertDetails = &ExpertDetails{
			FirstName:      r.ExpertDetails.FirstName,
			MiddleName:     r.ExpertDetails.MiddleName,
			LastName:       r.ExpertDetails.LastName,
			Email:          r.ExpertDetails.Email,
			Specialization: r.ExpertDetails.Specialization,
		}
	}

	return user
}

func (s *ctrlImpl) searchRsFromPb(rs *pb.SearchResponse) *SearchResponse {
	r := &SearchResponse{
		Index: int(rs.Paging.Index),
		Total: int(rs.Paging.Total),
		Users: []*User{},
	}

	for _, t := range rs.Users {
		r.Users = append(r.Users, s.fromPb(t))
	}

	return r
}
