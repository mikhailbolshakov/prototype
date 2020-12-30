package grpc

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/users/domain"
)

func (s *Server) fromPb(request *pb.CreateUserRequest) *domain.User {
	return &domain.User{
		Type:      request.Type,
		Username:  request.Username,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		Email:     request.Email,
	}
}

func (s *Server) fromDomain(user *domain.User) *pb.User {

	if user == nil {
		return nil
	}

	return &pb.User{
			Id:        user.Id,
			Type:      user.Type,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
			MMId:      user.MMUserId,
			MMChannelId: user.MMChannelId,
		}
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
