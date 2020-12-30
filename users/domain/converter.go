package domain

import (
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
)

func toDto(domain *User) *storage.User {
	return &storage.User{
		BaseDto:     kit.BaseDto{},
		Id:          domain.Id,
		Type:        domain.Type,
		Username:    domain.Username,
		FirstName:   domain.FirstName,
		LastName:    domain.LastName,
		Phone:       domain.Phone,
		Email:       domain.Email,
		MMUserId:    domain.MMUserId,
		MMChannelId: domain.MMChannelId,
	}
}

func fromDto(dto *storage.User) *User {

	if dto == nil {
		return nil
	}

	return &User{
		Id:          dto.Id,
		Type:        dto.Type,
		Username:    dto.Username,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Phone:       dto.Phone,
		Email:       dto.Email,
		MMUserId:    dto.MMUserId,
		MMChannelId: dto.MMChannelId,
	}
}

func toMM(domain *User) *mattermost.MMCreateUserRequest {
	return &mattermost.MMCreateUserRequest{
		Username: domain.Username,
		Email:    domain.Email,
	}
}

func criteriaToDto(c *SearchCriteria) *storage.SearchCriteria {
	if c == nil {
		return nil
	}

	return &storage.SearchCriteria{
		PagingRequest: c.PagingRequest,
		UserType:      c.UserType,
		Username:      c.Username,
		Email:         c.Email,
		Phone:         c.Phone,
		MMId:          c.MMId,
		MMChannelId:   c.MMChannelId,
	}
}

func searchRsFromDto(rs *storage.SearchResponse) *SearchResponse {
	if rs == nil {
		return nil
	}

	r := &SearchResponse{
		PagingResponse: rs.PagingResponse,
		Users: []*User{},
	}

	for _, t := range rs.Users {
		r.Users = append(r.Users, fromDto(t))
	}

	return r

}
