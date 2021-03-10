package grpc

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	"gitlab.medzdrav.ru/prototype/chat/logger"
	"gitlab.medzdrav.ru/prototype/chat/meta"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"gitlab.medzdrav.ru/prototype/proto/config"
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
	gs, err := kitGrpc.NewServer(meta.ServiceCode, logger.LF())
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterUsersServer(s.Srv, s)
	pb.RegisterChannelsServer(s.Srv, s)
	pb.RegisterPostsServer(s.Srv, s)

	return s
}

func (s *Server) Init(c *config.Config) error {
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

	rs, err := s.domain.CreateUser(ctx, &domain.CreateUserRequest{
		Username: rq.Username,
		Email:    rq.Email,
	})
	if err != nil {
		return nil, err
	}
	response := &pb.CreateUserResponse{ChatUserId: rs.Id}

	return response, nil
}

func (s *Server) SetStatus(ctx context.Context, rq *pb.SetStatusRequest) (*pb.SetStatusResponse, error) {

	domainRq := &domain.SetUserStatusRequest{
		ChatUserId: rq.UserStatus.ChatUserId,
		Status:     rq.UserStatus.Status,
		From:       s.toFromDomain(rq.From),
	}

	err := s.domain.SetStatus(ctx, domainRq)
	if err != nil {
		return nil, err
	}
	response := &pb.SetStatusResponse{}
	return response, nil
}

func (s *Server) CreateClientChannel(ctx context.Context, rq *pb.CreateClientChannelRequest) (*pb.CreateClientChannelResponse, error) {

	rs, err := s.domain.CreateClientChannel(ctx, &domain.CreateClientChannelRequest{
		ChatUserId:  rq.ChatUserId,
		DisplayName: rq.DisplayName,
		Name:        rq.Name,
		Subscribers: rq.Subscribers,
	})
	if err != nil {
		return nil, err
	}
	response := &pb.CreateClientChannelResponse{ChannelId: rs.ChannelId}

	return response, nil
}

func (s *Server) Subscribe(ctx context.Context, rq *pb.SubscribeRequest) (*pb.SubscribeResponse, error) {
	err := s.domain.SubscribeUser(ctx, &domain.SubscribeUserRequest{
		ChatUserId: rq.ChatUserId,
		ChannelId:  rq.ChannelId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SubscribeResponse{}, nil
}

func (s *Server) GetChannelsForUserAndMembers(ctx context.Context, rq *pb.GetChannelsForUserAndMembersRequest) (*pb.GetChannelsForUserAndMembersResponse, error) {

	channels, err := s.domain.GetChannelsForUserAndMembers(ctx, &domain.GetChannelsForUserAndMembersRequest{
		UserId:        rq.ChatUserId,
		MemberUserIds: rq.MemberChatUserIds,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetChannelsForUserAndMembersResponse{ChannelIds: channels}, nil

}

func (s *Server) GetUsersStatuses(ctx context.Context, rq *pb.GetUsersStatusesRequest) (*pb.GetUserStatusesResponse, error) {

	rs, err := s.domain.GetUsersStatuses(ctx, &domain.GetUsersStatusesRequest{ChatUserIds: rq.ChatUserIds})
	if err != nil {
		return nil, err
	}
	response := &pb.GetUserStatusesResponse{Statuses: []*pb.UserStatus{}}
	for _, s := range rs.Statuses {
		response.Statuses = append(response.Statuses, &pb.UserStatus{
			Status:     s.Status,
			ChatUserId: s.ChatUserId,
		})
	}
	return response, nil

}

func (s *Server) Post(ctx context.Context, rq *pb.PostRequest) (*pb.PostResponse, error) {

	var dPosts []*domain.Post
	for _, post := range rq.Posts {

		dPost := &domain.Post{
			Message:      post.Message,
			ToChatUserId: post.ToChatUserId,
			ChannelId:    post.ChannelId,
			Ephemeral:    post.Ephemeral,
			From:         s.toFromDomain(post.From),
			Attachments:  []*domain.PostAttachment{},
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

	err := s.domain.Posts(ctx, dPosts)
	if err != nil {
		return nil, err
	}

	return &pb.PostResponse{}, nil

}

func (s *Server) AskBot(ctx context.Context, rq *pb.AskBotRequest) (*pb.AskBotResponse, error) {
	rs, err := s.domain.AskBot(ctx, &domain.AskBotRequest{
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
	err := s.domain.DeleteUser(ctx, rq.ChatUserId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{}, nil
}

func (s *Server) Login(ctx context.Context, rq *pb.LoginRequest) (*pb.LoginResponse, error) {
	rs, err := s.domain.Login(ctx, &domain.LoginRequest{
		UserId:     rq.UserId,
		Username:   rq.Username,
		ChatUserId: rq.ChatUserId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{ChatSessionId: rs.ChatSessionId}, nil
}
func (s *Server) Logout(ctx context.Context, rq *pb.LogoutRequest) (*pb.LogoutResponse, error) {

	err := s.domain.Logout(ctx, &domain.LogoutRequest{
		ChatUserId: rq.ChatUserId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResponse{}, nil
}
