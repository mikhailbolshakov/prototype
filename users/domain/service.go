package domain

import (
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/users/repository"
)

type UserService interface {
	Create(user *User) (*User, error)
	GetByUsername(username string) *User
}

type UserServiceImpl struct {
	storage    repository.UserStorage
	muttermost repository.Mattermost
}

func NewUserService() UserService {
	return &UserServiceImpl{
		storage:    repository.NewStorage(),
		muttermost: repository.NewMattermost(),
	}
}

func (u *UserServiceImpl) Create(user *User) (*User, error) {

	// create a new user in MM
	mmRq, err := u.toMM(user)
	if err != nil {
		return nil, err
	}
	mmRs, err := u.muttermost.CreateUser(mmRq)
	if err != nil {
		return nil, err
	}

	user.MMUserId = mmRs.Id
	user.Id = kit.NewId()

	// save to storage
	dto, err := u.toDto(user)
	if err != nil {
		return nil, err
	}

	dto, err = u.storage.CreateUser(dto)
	if err != nil {
		return nil, err
	}

	return u.fromDto(dto)

}

func (u *UserServiceImpl) GetByUsername(username string) *User {
	r, _ := u.fromDto(u.storage.GetByUsername(username))
	return r
}
