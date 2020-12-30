package mattermost

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/mm"
)

type MMNewPostMessageHandler func(post *MMPost)

type Service interface {
	CreateUser(user *MMCreateUserRequest) (*MMCreateUserResponse, error)
	CreateClientChannel(rq *MMCreateClientChannelRequest) (*MMCreateClientChannelResponse, error)
	SetNewPostMessageHandler(handler MMNewPostMessageHandler)
}

type MMCreateUserRequest struct {
	Username string
	Email string
}

type MMCreateUserResponse struct {
	Id string
}

type MMCreateClientChannelRequest struct {
	ClientUserId string
}

type MMCreateClientChannelResponse struct {
	ChannelId string
}

type MMPost struct {
	Id        string `json:"id"`
	CreateAt  int64  `json:"createAt"`
	UserId    string `json:"userId"`
	ChannelId string `json:"channelId"`
	Message   string `json:"message"`
}

type serviceImpl struct {
	client *mm.Client
	queue  queue.Queue
	newPostsHandler MMNewPostMessageHandler
}

func newImpl() *serviceImpl {
	m := &serviceImpl{}
	return m
}

func (s *serviceImpl) CreateUser(rq *MMCreateUserRequest) (*MMCreateUserResponse, error) {

	rs, err := s.client.CreateUser(&mm.CreateUserRequest{
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

func (s *serviceImpl) CreateClientChannel(rq *MMCreateClientChannelRequest) (*MMCreateClientChannelResponse, error) {
	rs, err := s.client.CreateClientChannel(&mm.CreateClientChannelRequest{
		ClientUserId: rq.ClientUserId,
		TeamName:     "rgs",
	})
	if err != nil {
		return nil, err
	}

	return &MMCreateClientChannelResponse{ChannelId: rs.ChannelId}, nil
}

func (s *serviceImpl) SetNewPostMessageHandler(handler MMNewPostMessageHandler) {
	s.newPostsHandler = handler
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
				var post mm.Post
				_ = json.Unmarshal(msg, &post)
				s.newPostsHandler(&MMPost{
					Id:        post.Id,
					CreateAt:  post.CreateAt,
					UserId:    post.UserId,
					ChannelId: post.ChannelId,
					Message:   post.Message,
				})
			}
		}

	}()

	return nil
}
