package impl

import (
	"context"
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

func (s *serviceImpl) GetChannelsForUserAndMembers(ctx context.Context, rq *domain.GetChannelsForUserAndMembersRequest) ([]string, error) {
	return s.mmService.GetChannelsForUserAndMembers(ctx, rq)
}

func (s *serviceImpl) GetUsersStatuses(ctx context.Context, rq *domain.GetUsersStatusesRequest) (*domain.GetUsersStatusesResponse, error) {
	return s.mmService.GetUserStatuses(ctx, rq)
}

func (s *serviceImpl) CreateUser(ctx context.Context, rq *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {

	userId, err := s.mmService.CreateUser(ctx, rq)

	if err != nil {
		return nil, err
	}

	return &domain.CreateUserResponse{Id: userId}, nil
}

func (s *serviceImpl) CreateClientChannel(ctx context.Context, rq *domain.CreateClientChannelRequest) (*domain.CreateClientChannelResponse, error) {

	channelId, err := s.mmService.CreateClientChannel(ctx, rq)
	if err != nil {
		return nil, err
	}

	if rq.Subscribers != nil && len(rq.Subscribers) > 0 {
		for _, sbUserId := range rq.Subscribers {
			err = s.mmService.SubscribeUser(ctx, sbUserId, channelId)
			if err != nil {
				return nil, err
			}
		}
	}

	return &domain.CreateClientChannelResponse{ChannelId: channelId}, nil
}

func (s *serviceImpl) SubscribeUser(ctx context.Context, rq *domain.SubscribeUserRequest) error {
	return s.mmService.SubscribeUser(ctx, rq.UserId, rq.ChannelId)
}

func (s *serviceImpl) DeleteUser(ctx context.Context, userId string) error {
	return s.mmService.DeleteUser(ctx, userId)
}

func (s *serviceImpl) Posts(ctx context.Context, posts []*domain.Post) error {

	var err error
	for _, post := range posts {

		if post.Ephemeral && post.ToUserId == "" {
			return fmt.Errorf("recipient user id must be specified for an ephemeral post")
		}

		if post.PredefinedPost != nil && post.PredefinedPost.Code != "" {
			post, err = s.predefinedPost(ctx, post)
			if err != nil {
				return err
			}
		}

		if err := s.mmService.Post(ctx, post.ChannelId, post.Message, post.ToUserId, post.Ephemeral, post.FromBot, post.Attachments); err != nil {
			return err
		}

	}

	return nil
}