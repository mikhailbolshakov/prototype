package impl

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/chat/domain"
)

type serviceImpl struct {
	mmService    domain.MattermostService
}

func NewService(mmService domain.MattermostService) domain.Service {

	s := &serviceImpl{
		mmService:    mmService,
	}

	return s
}

func (s *serviceImpl) GetChannelsForUserAndMembers(rq *domain.GetChannelsForUserAndMembersRequest) ([]string, error) {
	return s.mmService.GetChannelsForUserAndMembers(rq)
}

func (s *serviceImpl) GetUsersStatuses(rq *domain.GetUsersStatusesRequest) (*domain.GetUsersStatusesResponse, error) {
	return s.mmService.GetUserStatuses(rq)
}

func (s *serviceImpl) CreateUser(rq *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {

	userId, err := s.mmService.CreateUser(rq)

	if err != nil {
		return nil, err
	}

	return &domain.CreateUserResponse{Id: userId}, nil
}

func (s *serviceImpl) CreateClientChannel(rq *domain.CreateClientChannelRequest) (*domain.CreateClientChannelResponse, error) {

	channelId, err := s.mmService.CreateClientChannel(rq)
	if err != nil {
		return nil, err
	}

	if rq.Subscribers != nil && len(rq.Subscribers) > 0 {
		for _, sbUserId := range rq.Subscribers {
			err = s.mmService.SubscribeUser(sbUserId, channelId)
			if err != nil {
				return nil, err
			}
		}
	}

	return &domain.CreateClientChannelResponse{ChannelId: channelId}, nil
}

func (s *serviceImpl) SubscribeUser(rq *domain.SubscribeUserRequest) error {
	return s.mmService.SubscribeUser(rq.UserId, rq.ChannelId)
}

func (s *serviceImpl) DeleteUser(userId string) error {
	return s.mmService.DeleteUser(userId)
}

func (s *serviceImpl) Posts(posts []*domain.Post) error {

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

		if err := s.mmService.Post(post.ChannelId, post.Message, post.ToUserId, post.Ephemeral, post.FromBot, post.Attachments); err != nil {
			return err
		}

	}

	return nil
}