package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/users/repository"
)

func (u *UserServiceImpl) toDto(domain *User) (*repository.User, error) {
	return &repository.User{
		BaseDto:   storage.BaseDto{},
		Id:        domain.Id,
		Type:      domain.Type,
		Username:  domain.Username,
		FirstName: domain.FirstName,
		LastName:  domain.LastName,
		Phone:     domain.Phone,
		Email:     domain.Email,
		MMUserId:  domain.MMUserId,
	}, nil
}

func (u *UserServiceImpl) fromDto(dto *repository.User) (*User, error) {

	if dto == nil {
		return nil, nil
	}

	return &User{
		Id:        dto.Id,
		Type:      dto.Type,
		Username:  dto.Username,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Phone:     dto.Phone,
		Email:     dto.Email,
		MMUserId:  dto.MMUserId,
	}, nil
}

func (u *UserServiceImpl) toMM(domain *User) (*repository.MMCreateUserRequest, error) {
	return &repository.MMCreateUserRequest{
		Username: domain.Username,
		Email:    domain.Email,
	}, nil
}
