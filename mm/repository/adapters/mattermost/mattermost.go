package mattermost

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/queue"
)

type NewPostMessageHandler func(post *Post)

type Service interface {
	CreateUser(user *CreateUserRequest) (*CreateUserResponse, error)
	CreateClientChannel(rq *CreateClientChannelRequest) (*CreateClientChannelResponse, error)
	SubscribeUser(rq *SubscribeUserRequest) (*SubscribeUserResponse, error)
	SetNewPostMessageHandler(handler NewPostMessageHandler)
	CreateEphemeralPost(p *EphemeralPost) error
	CreatePost(p *Post) error
}

type serviceImpl struct {
	client          *Client
	queue           queue.Queue
	newPostsHandler NewPostMessageHandler
}

func newImpl(queue queue.Queue) *serviceImpl {
	m := &serviceImpl{
		queue: queue,
	}
	return m
}

func (s *serviceImpl) CreateUser(rq *CreateUserRequest) (*CreateUserResponse, error) {

	if rq.TeamName == "" {
		rq.TeamName = "rgs"
	}

	if rq.Password == "" {
		rq.Password = "12345"
	}

	rs, err := s.client.createUser(rq)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (s *serviceImpl) CreateClientChannel(rq *CreateClientChannelRequest) (*CreateClientChannelResponse, error) {

	if rq.TeamName == "" {
		rq.TeamName = "rgs"
	}

	rs, err := s.client.createClientChannel(rq)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (s *serviceImpl) SetNewPostMessageHandler(handler NewPostMessageHandler) {
	s.newPostsHandler = handler
}

func (s *serviceImpl) SubscribeUser(rq *SubscribeUserRequest) (*SubscribeUserResponse, error) {
	err := s.client.subscribeUser(rq.ChannelId, rq.UserId)
	if err != nil {
		return nil, err
	}

	return &SubscribeUserResponse{}, nil
}

func (s *serviceImpl) CreateEphemeralPost(p *EphemeralPost) error {
	return s.client.createEphemeralPost(p.ChannelId, p.UserId, p.Message, p.Attachments)
}

func (s *serviceImpl) CreatePost(p *Post) error {
	return s.client.createPost(p.ChannelId, p.Message, p.Attachments)
}

func (s *serviceImpl) listenNewPosts() error {

	receiver := make(chan []byte)
	err := s.queue.Subscribe("mm.posts", receiver)
	if err != nil {
		return err
	}

	go func() {

		for {
			select {
			case msg := <-receiver:
				post := &Post{}
				_ = json.Unmarshal(msg, post)
				s.newPostsHandler(post)
			}
		}

	}()

	return nil
}