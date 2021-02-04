package domain

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/mattermost"
)

type Service interface {
	GetUsersStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateUser(rq *CreateUserRequest) (*CreateUserResponse, error)
	CreateClientChannel(rq *CreateClientChannelRequest) (*CreateClientChannelResponse, error)
	GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	SubscribeUser(rq *SubscribeUserRequest) error
	DeleteUser(userId string) error
	AskBot(request *AskBotRequest) (*AskBotResponse, error)
	Posts(posts []*Post) error
}

type serviceImpl struct {
	mmService    mattermost.Service
}

func NewService(mmService mattermost.Service) Service {

	s := &serviceImpl{
		mmService:    mmService,
	}

	return s
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

func (s *serviceImpl) Posts(posts []*Post) error {

	var err error
	for _, post := range posts {

		if post.Ephemeral && post.ToUserId == "" {
			return fmt.Errorf("recipient user id must be specified for an ephemeral post")
		}

		if post.PredefinedPost != nil && post.PredefinedPost.Code != "" {
			post, err = s.predefinedPost(post)
			if err != nil {
				return err
			}
		}

		if err := s.mmService.Post(post.ChannelId, post.Message, post.ToUserId, post.Ephemeral, post.FromBot, s.convertAttachments(post.Attachments)); err != nil {
			return err
		}

	}

	return nil
}