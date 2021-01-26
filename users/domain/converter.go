package domain

import (
	"encoding/json"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
)

func toDto(domain *User) *storage.User {

	dto := &storage.User{
		BaseDto:  kit.BaseDto{},
		Id:       domain.Id,
		Type:     domain.Type,
		Status:   domain.Status,
		Username: domain.Username,
		MMUserId: domain.MMUserId,
		KKUserId: domain.KKUserId,
	}

	var detailsBytes []byte
	switch domain.Type {
	case USER_TYPE_CLIENT:
		detailsBytes, _ = json.Marshal(domain.ClientDetails)
	case USER_TYPE_CONSULTANT:
		detailsBytes, _ = json.Marshal(domain.ConsultantDetails)
	case USER_TYPE_EXPERT:
		detailsBytes, _ = json.Marshal(domain.ExpertDetails)
	}

	dto.Details = string(detailsBytes)

	return dto

}

func fromDto(dto *storage.User) *User {

	if dto == nil {
		return nil
	}

	domain := &User{
		Id:         dto.Id,
		Type:       dto.Type,
		Username:   dto.Username,
		Status:     dto.Status,
		MMUserId:   dto.MMUserId,
		KKUserId:   dto.KKUserId,
		ModifiedAt: dto.UpdatedAt,
		DeletedAt:  dto.DeletedAt,
	}

	switch dto.Type {
	case USER_TYPE_CLIENT:
		cd := &ClientDetails{}
		_ = json.Unmarshal([]byte(dto.Details), cd)
		domain.ClientDetails = cd
	case USER_TYPE_CONSULTANT:
		cd := &ConsultantDetails{}
		_ = json.Unmarshal([]byte(dto.Details), cd)
		domain.ConsultantDetails = cd
	case USER_TYPE_EXPERT:
		ed := &ExpertDetails{}
		_ = json.Unmarshal([]byte(dto.Details), &ed)
		domain.ExpertDetails = ed
	}

	return domain

}

func criteriaToDto(c *SearchCriteria) *storage.SearchCriteria {
	if c == nil {
		return nil
	}

	return &storage.SearchCriteria{
		PagingRequest: c.PagingRequest,
		UserType:      c.UserType,
		Username:      c.Username,
		Status:        c.Status,
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
		Users:          []*User{},
	}

	for _, t := range rs.Users {
		r.Users = append(r.Users, fromDto(t))
	}

	return r

}
