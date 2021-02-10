package domain

import (
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
)

type UserStorage interface {
	CreateUser(u *User) (*User, error)
	GetByUsername(username string) *User
	GetByMMId(mmId string) *User
	Get(id string) *User
	Search(cr *SearchCriteria) (*SearchResponse, error)
	UpdateStatus(userId, status string, isDeleted bool) (*User, error)
	UpdateDetails(userId string, details string) (*User, error)
	UpdateMMId(userId, mmId string) (*User, error)
	UpdateKKId(userId, kkId string) (*User, error)
}

type ChatService interface {
	CreateUser(rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	CreateClientChannel(rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error)
	GetUsersStatuses(rq *pb.GetUsersStatusesRequest) (*pb.GetUserStatusesResponse, error)
}

type ConfigService interface {
	Get() (*kit.Config, error)
}