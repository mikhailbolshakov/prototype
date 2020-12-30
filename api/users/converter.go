package users

import (
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

func (c *controller) toPb(request *CreateUserRequest) *pb.CreateUserRequest {
	return &pb.CreateUserRequest{
		Username:  request.Username,
		Type:      request.Type,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		Email:     request.Email,
	}
}

func (s *controller) fromPb(response *pb.User) *User {
	return &User{
			Type:        response.Type,
			Username:    response.Username,
			FirstName:   response.FirstName,
			LastName:    response.LastName,
			Phone:       response.Phone,
			Email:       response.Email,
			Id:          response.Id,
			MMId:        response.MMId,
			MMChannelId: response.MMChannelId,
		}
}

func (s *controller) searchRsFromPb(rs *pb.SearchResponse) *SearchResponse {
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
