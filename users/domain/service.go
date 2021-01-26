package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
	"strings"
	"time"
)

type UserService interface {
	Create(user *User) (*User, error)
	GetByUsername(username string) *User
	GetByMMId(mmId string) *User
	Get(id string) *User
	Activate(userId string) (*User, error)
	Delete(userId string) (*User, error)
	SetClientDetails(userId string, details *ClientDetails) (*User, error)
	SetMMUserId(userId, mmId string) (*User, error)
	SetKKUserId(userId, kkId string) (*User, error)
}

type userServiceImpl struct {
	common.BaseService
	storage    storage.UserStorage
	mattermost mattermost.Service
}

func NewUserService(storage storage.UserStorage, mmService mattermost.Service, queue queue.Queue) UserService {

	s := &userServiceImpl{
		storage:    storage,
		mattermost: mmService,
	}
	s.BaseService = common.BaseService{Queue: queue}

	return s
}

func (u *userServiceImpl) newClient(user *User) (*User, error) {

	if user.ClientDetails == nil {
		return nil, fmt.Errorf("details isn't populated properly")
	}

	if user.ClientDetails.Phone == "" {
		return nil, fmt.Errorf("phone is empty")
	}

	var sex = map[string]bool{"M": true, "F": true}

	if _, ok := sex[user.ClientDetails.Sex]; !ok {
		return nil, fmt.Errorf("sex is incorrect")
	}

	if user.ClientDetails.FirstName == "" ||
		user.ClientDetails.LastName == "" ||
		user.ClientDetails.BirthDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) ||
		user.ClientDetails.BirthDate.After(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return nil, fmt.Errorf("pesonal data is incorrect")
	}

	if user.ClientDetails.PersonalAgreement == nil {
		user.ClientDetails.PersonalAgreement = &PersonalAgreement{}
	}

	user.Username = user.ClientDetails.Phone

	return user, nil

}

func (u *userServiceImpl) newConsultant(user *User) (*User, error) {

	if user.ConsultantDetails == nil {
		return nil, fmt.Errorf("details isn't populated properly")
	}

	if user.ConsultantDetails.Email == "" {
		return nil, fmt.Errorf("email is empty")
	}

	if user.ConsultantDetails.FirstName == "" ||
		user.ConsultantDetails.LastName == "" {
		return nil, fmt.Errorf("pesonal data is incorrect")
	}

	user.Username = strings.Split(user.ConsultantDetails.Email, "@")[0]

	return user, nil

}

func (u *userServiceImpl) newExpert(user *User) (*User, error) {

	if user.ExpertDetails == nil {
		return nil, fmt.Errorf("details isn't populated properly")
	}

	if user.ExpertDetails.Email == "" {
		return nil, fmt.Errorf("email is empty")
	}

	if user.ExpertDetails.FirstName == "" ||
		user.ExpertDetails.LastName == "" {
		return nil, fmt.Errorf("pesonal data is incorrect")
	}

	user.Username = strings.Split(user.ExpertDetails.Email, "@")[0]

	return user, nil

}

func (u *userServiceImpl) Create(user *User) (*User, error) {

	user.Id = kit.NewId()
	user.Status = USER_STATUS_DRAFT

	var err error
	switch user.Type {
	case USER_TYPE_CLIENT:
		user, err = u.newClient(user)
	case USER_TYPE_CONSULTANT:
		user, err = u.newConsultant(user)
	case USER_TYPE_EXPERT:
		user, err = u.newExpert(user)
	case USER_TYPE_SUPERVISOR:
		return nil, errors.New("not implemented")
	default:
		return nil, fmt.Errorf("not supported user type %s", user.Type)
	}
	if err != nil {
		return nil, err
	}

	// check username uniqueness
	if usr := u.storage.GetByUsername(user.Username); usr != nil && usr.Id != "" {
		return nil, fmt.Errorf("username %s already exists", user.Username)
	}

	// save to storage
	dto, err := u.storage.CreateUser(toDto(user))
	if err != nil {
		return nil, err
	}

	user = fromDto(dto)

	u.Publish(user, "users.draft-created")

	return user, nil

}

func (u *userServiceImpl) GetByUsername(username string) *User {
	return fromDto(u.storage.GetByUsername(username))
}

func (u *userServiceImpl) GetByMMId(mmId string) *User {
	return fromDto(u.storage.GetByMMId(mmId))
}

func (u *userServiceImpl) Get(id string) *User {
	return fromDto(u.storage.Get(id))
}

func (u *userServiceImpl) Activate(userId string) (*User, error) {

	dto := u.storage.Get(userId)
	if dto == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	dto, err := u.storage.UpdateStatus(userId, USER_STATUS_ACTIVE, false)
	if err != nil {
		return nil, err
	}

	return fromDto(dto), nil
}

func (u *userServiceImpl) Delete(userId string) (*User, error) {

	dto := u.storage.Get(userId)
	if dto == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	dto, err := u.storage.UpdateStatus(userId, USER_STATUS_DELETED, true)
	if err != nil {
		return nil, err
	}

	return fromDto(dto), nil
}

func (u *userServiceImpl) SetClientDetails(userId string, details *ClientDetails) (*User, error) {

	dto := u.storage.Get(userId)
	if dto == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	if dto.Type != USER_TYPE_CLIENT {
		return nil, fmt.Errorf("user withid %s isn't a client", userId)
	}

	detB, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	dto, err = u.storage.UpdateDetails(userId, string(detB))
	if err != nil {
		return nil, err
	}

	return fromDto(dto), nil
}

func (u *userServiceImpl) SetMMUserId(userId, mmId string) (*User, error) {

	dto := u.storage.Get(userId)
	if dto == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	dto, err := u.storage.UpdateMMId(userId, mmId)
	if err != nil {
		return nil, err
	}

	return fromDto(dto), nil
}

func (u *userServiceImpl) SetKKUserId(userId, kkId string) (*User, error) {

	dto := u.storage.Get(userId)
	if dto == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	dto, err := u.storage.UpdateKKId(userId, kkId)
	if err != nil {
		return nil, err
	}

	return fromDto(dto), nil
}