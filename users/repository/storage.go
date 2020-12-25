package repository

import "time"

type UserStorage interface {
	CreateUser(u *User) (*User, error)
	GetByUsername(username string) *User
}

type UserStorageImpl struct {}

func NewStorage() UserStorage {
	return &UserStorageImpl{}
}

func (s *UserStorageImpl) CreateUser(user *User) (*User, error) {

	t := time.Now()
	user.CreatedAt, user.UpdatedAt = t, t

	result := storage.Instance.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s *UserStorageImpl) GetByUsername(username string) *User {

	user := &User{}

	storage.Instance.Where("username = ?", username).First(&user)

	return user
}
