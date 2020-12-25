package repository

import "gitlab.medzdrav.ru/prototype/mm"

type Mattermost interface {
	CreateUser(user *MMCreateUserRequest) (*MMCreateUserResponse, error)
}

type MattermostImpl struct {
}

func NewMattermost() Mattermost {
	return &MattermostImpl{}
}

func (s *MattermostImpl) CreateUser(rq *MMCreateUserRequest) (*MMCreateUserResponse, error) {

	rs, err := mmClient.CreateUser(&mm.CreateUserRequest{
		TeamName: "rgs",
		Username: rq.Username,
		Password: "12345",
		Email:    rq.Email,
	})
	if err != nil {
		return nil, err
	}

	return &MMCreateUserResponse{
		Id: rs.Id,
	}, nil
}