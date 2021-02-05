package grpc

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"log"
)

type Server struct {
	host, port string
	*kitGrpc.Server
	domain domain.Service
	pb.UnimplementedUsersServer
	pb.UnimplementedChannelsServer
	pb.UnimplementedPostsServer
}

func New(domain domain.Service) *Server {

	s := &Server{domain: domain}

	// grpc server
	gs, err := kitGrpc.NewServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterUsersServer(s.Srv, s)
	pb.RegisterChannelsServer(s.Srv, s)
	pb.RegisterPostsServer(s.Srv, s)

	return s
}

func  (s *Server) Init(c *kitConfig.Config) error {
	usersCfg := c.Services["chat"]
	s.host = usersCfg.Grpc.Host
	s.port = usersCfg.Grpc.Port
	return nil
}

func (s *Server) ListenAsync() {

	go func() {
		err := s.Server.Listen(s.host, s.port)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) CreateUser(ctx context.Context, rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	rs, err := s.domain.CreateUser(&domain.CreateUserRequest{
		Username: rq.Username,
		Email:    rq.Email,
	})
	if err != nil {
		return nil, err
	}
	response := &pb.CreateUserResponse{Id: rs.Id}

	return response, nil
}

func (s *Server) CreateClientChannel(ctx context.Context, rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error) {

	rs, err := s.domain.CreateClientChannel(&domain.CreateClientChannelRequest{
		ClientUserId: rq.ClientUserId,
		DisplayName:  rq.DisplayName,
		Name:         rq.Name,
		Subscribers:  rq.Subscribers,
	})
	if err != nil {
		return nil, err
	}
	response := &pb.CreateClientChannelResponse{ChannelId: rs.ChannelId}

	return response, nil
}

func (s *Server) Subscribe(ctx context.Context, rq *pb.SubscribeRequest) (*pb.SubscribeResponse, error) {
	err := s.domain.SubscribeUser(&domain.SubscribeUserRequest{
		UserId:    rq.UserId,
		ChannelId: rq.ChannelId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SubscribeResponse{}, nil
}

func (s *Server) GetChannelsForUserAndMembers(ctx context.Context, rq *pb.GetChannelsForUserAndMembersRequest) (*pb.GetChannelsForUserAndMembersResponse, error) {

	channels, err := s.domain.GetChannelsForUserAndMembers(&domain.GetChannelsForUserAndMembersRequest{
		UserId:        rq.UserId,
		TeamName:      rq.TeamName,
		MemberUserIds: rq.MemberUserIds,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetChannelsForUserAndMembersResponse{ChannelIds: channels}, nil

}

func (s *Server) GetUsersStatuses(ctx context.Context, rq *pb.GetUsersStatusesRequest) (*pb.GetUserStatusesResponse, error) {

	rs, err := s.domain.GetUsersStatuses(&domain.GetUsersStatusesRequest{UserIds: rq.MMUserIds})
	if err != nil {
		return nil, err
	}
	response := &pb.GetUserStatusesResponse{Statuses: []*pb.UserStatus{}}
	for _, s := range rs.Statuses {
		response.Statuses = append(response.Statuses, &pb.UserStatus{
			Status:   s.Status,
			MMUserId: s.UserId,
		})
	}
	return response, nil

}

func (s *Server) Post(ctx context.Context, rq *pb.PostRequest) (*pb.PostResponse, error) {

	var dPosts []*domain.Post
	for _, post := range rq.Posts {

		dPost := &domain.Post{
			Message:     post.Message,
			ToUserId:    post.ToUserId,
			ChannelId:   post.ChannelId,
			Ephemeral:   post.Ephemeral,
			FromBot:     post.FromBot,
			Attachments: []*domain.PostAttachment{},
		}

		if post.PredefinedPost != nil && post.PredefinedPost.Code != "" {
			dPost.PredefinedPost = &domain.PredefinedPost{
				Code: post.PredefinedPost.Code,
			}

			if post.PredefinedPost.Params != nil {
				var params map[string]interface{}
				if err := json.Unmarshal(post.PredefinedPost.Params, &params); err != nil {
					return nil, err
				}
				dPost.PredefinedPost.Params = params
			}

		}

		if post.Attachments != nil && len(post.Attachments) > 0 {

			for _, att := range post.Attachments {

				dAtt := &domain.PostAttachment{
					Fallback:   att.Fallback,
					Color:      att.Color,
					Pretext:    att.Pretext,
					AuthorName: att.AuthorName,
					AuthorLink: att.AuthorLink,
					AuthorIcon: att.AuthorIcon,
					Title:      att.Title,
					TitleLink:  att.TitleLink,
					Text:       att.Text,
					ImageURL:   att.ImageURL,
					ThumbURL:   att.ThumbURL,
					Footer:     att.Footer,
					FooterIcon: att.FooterIcon,
				}

				dPost.Attachments = append(dPost.Attachments, dAtt)
			}

		}

		dPosts = append(dPosts, dPost)

	}

	err := s.domain.Posts(dPosts)
	if err != nil {
		return nil, err
	}

	return &pb.PostResponse{}, nil

}

func (s *Server) AskBot(ctx context.Context, rq *pb.AskBotRequest) (*pb.AskBotResponse, error) {
	rs, err := s.domain.AskBot(&domain.AskBotRequest{
		From:    rq.From,
		Message: rq.Message,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AskBotResponse{
		Found:  rs.Found,
		Answer: rs.Answer,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, rq *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.domain.DeleteUser(rq.MMUserId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{}, nil
}
