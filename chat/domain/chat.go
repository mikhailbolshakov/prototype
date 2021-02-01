package domain

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/bp"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/users"
)

type Service interface {
	GetUsersStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateUser(rq *CreateUserRequest) (*CreateUserResponse, error)
	CreateClientChannel(rq *CreateClientChannelRequest) (*CreateClientChannelResponse, error)
	GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	SendTriggerPost(rq *SendTriggerPostRequest) error
	SendPostFromBot(rq *SendPostRequest) error
	SubscribeUser(rq *SubscribeUserRequest) error
	DeleteUser(userId string) error
	AskBot(request *AskBotRequest) (*AskBotResponse, error)
}

type serviceImpl struct {
	mmService    mattermost.Service
	usersService users.Service
	tasksService tasks.Service
	bpService    bp.Service
}

func NewService(mmService mattermost.Service, usersService users.Service, tasksService tasks.Service, bpService bp.Service) Service {

	s := &serviceImpl{
		mmService:    mmService,
		usersService: usersService,
		tasksService: tasksService,
		bpService:    bpService,
	}

	return s
}

func (s *serviceImpl) SendPostFromBot(rq *SendPostRequest) error {

	post := &mattermost.Post{
		ChannelId:   rq.ChannelId,
		Message:     rq.Message,
	}

	if rq.Ephemeral {
		if rq.UserId == "" {
			return fmt.Errorf("user isn't specified for ephemeral post")
		}
		if err := s.mmService.EphemeralPostFromBot(&mattermost.EphemeralPost{
			Post:   post,
			UserId: rq.UserId,
		}); err != nil {
			return err
		}
	} else {
		if err := s.mmService.PostFromBot(post); err != nil {
			return err
		}
	}

	return nil
}

func (s *serviceImpl) GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error) {
	return s.mmService.GetChannelsForUserAndMembers(&mattermost.GetChannelsForUserAndMembersRequest{
		UserId:        rq.UserId,
		TeamName:      rq.TeamName,
		MemberUserIds: rq.MemberUserIds,
	})
}

func (s *serviceImpl) GetUsersStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error) {

	rs, err := s.mmService.GetUserStatuses(&mattermost.GetUsersStatusesRequest{UserIds: rq.UserIds})
	if err != nil {
		return nil, err
	}

	response := &GetUsersStatusesResponse{Statuses: []*UserStatus{}}
	for _, s := range rs.Statuses {
		response.Statuses = append(response.Statuses, &UserStatus{
			UserId: s.UserId,
			Status: s.Status,
		})
	}

	return response, nil

}

func (s *serviceImpl) CreateUser(rq *CreateUserRequest) (*CreateUserResponse, error) {

	rs, err := s.mmService.CreateUser(&mattermost.CreateUserRequest{
		TeamName: rq.TeamName,
		Username: rq.Username,
		Password: rq.Password,
		Email:    rq.Email,
	})

	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{Id: rs.Id}, nil
}

func (s *serviceImpl) CreateClientChannel(rq *CreateClientChannelRequest) (*CreateClientChannelResponse, error) {

	rs, err := s.mmService.CreateClientChannel(&mattermost.CreateClientChannelRequest{
		ClientUserId: rq.ClientUserId,
		TeamName:     rq.TeamName,
		DisplayName:  rq.DisplayName,
		Name:         rq.Name,
	})
	if err != nil {
		return nil, err
	}

	if rq.Subscribers != nil && len(rq.Subscribers) > 0 {
		for _, sbUserId := range rq.Subscribers {
			_, err = s.mmService.SubscribeUser(&mattermost.SubscribeUserRequest{
				UserId:    sbUserId,
				ChannelId: rs.ChannelId,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return &CreateClientChannelResponse{ChannelId: rs.ChannelId}, nil
}

func (s *serviceImpl) SubscribeUser(rq *SubscribeUserRequest) error {
	_, err := s.mmService.SubscribeUser(&mattermost.SubscribeUserRequest{
		UserId:    rq.UserId,
		ChannelId: rq.ChannelId,
	})
	return err
}

func (s *serviceImpl) DeleteUser(userId string) error {
	return s.mmService.DeleteUser(userId)
}
