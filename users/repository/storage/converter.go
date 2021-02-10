package storage

import (
	"encoding/json"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/users/domain"
)

func (s *storageImpl) toUserDto(u *domain.User) *user {

	if u == nil {
		return nil
	}

	dto := &user{
		BaseDto:  kit.BaseDto{},
		Id:       u.Id,
		Type:     u.Type,
		Status:   u.Status,
		Username: u.Username,
		MMUserId: u.MMUserId,
		KKUserId: u.KKUserId,
	}

	var detailsBytes []byte
	switch u.Type {
	case domain.USER_TYPE_CLIENT:
		detailsBytes, _ = json.Marshal(u.ClientDetails)
	case domain.USER_TYPE_CONSULTANT:
		detailsBytes, _ = json.Marshal(u.ConsultantDetails)
	case domain.USER_TYPE_EXPERT:
		detailsBytes, _ = json.Marshal(u.ExpertDetails)
	}

	dto.Details = string(detailsBytes)

	dto.Groups = []*userGroup{}
	for _, g := range u.Groups {
		dto.Groups = append(dto.Groups, &userGroup{
			UserId: dto.Id,
			Type:   dto.Type,
			Group:  g,
		})
	}

	return dto

}

func (s *storageImpl) toUserDomain(dto *user) *domain.User {

	if dto == nil {
		return nil
	}

	d := &domain.User{
		Id:         dto.Id,
		Type:       dto.Type,
		Username:   dto.Username,
		Status:     dto.Status,
		MMUserId:   dto.MMUserId,
		KKUserId:   dto.KKUserId,
		ModifiedAt: dto.UpdatedAt,
		DeletedAt:  dto.DeletedAt,
		Groups:     []string{},
	}

	switch dto.Type {
	case domain.USER_TYPE_CLIENT:
		cd := &domain.ClientDetails{}
		_ = json.Unmarshal([]byte(dto.Details), cd)
		d.ClientDetails = cd
	case domain.USER_TYPE_CONSULTANT:
		cd := &domain.ConsultantDetails{}
		_ = json.Unmarshal([]byte(dto.Details), cd)
		d.ConsultantDetails = cd
	case domain.USER_TYPE_EXPERT:
		ed := &domain.ExpertDetails{}
		_ = json.Unmarshal([]byte(dto.Details), &ed)
		d.ExpertDetails = ed
	}
	for _, g := range dto.Groups {
		d.Groups = append(d.Groups, g.Group)
	}

	return d

}

func (s *storageImpl) toUsersDomain(dtos []*user) []*domain.User {
	var res []*domain.User
	for _, d := range dtos {
		res = append(res, s.toUserDomain(d))
	}
	return res
}

