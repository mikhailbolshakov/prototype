package grpc

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/chat/domain"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"log"
)

type Server struct {
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

func (s *Server) ListenAsync() {

	go func() {
		err := s.Server.Listen("localhost", "50053")
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

func (s *Server) SendTriggerPost(ctx context.Context, rq *pb.SendTriggerPostRequest) (*pb.SendTriggerPostResponse, error) {

	var params map[string]interface{}
	if rq.Params != nil {
		if err := json.Unmarshal(rq.Params, &params); err != nil {
			return nil, err
		}
	}

	domainRq := &domain.SendTriggerPostRequest{
		TriggerPostCode: rq.PostCode,
		UserId:          rq.UserId,
		ChannelId:       rq.ChannelId,
		Params:          params,
	}
	err := s.domain.SendTriggerPost(domainRq)
	if err != nil {
		return nil, err
	}

	return &pb.SendTriggerPostResponse{}, nil
}

func (s *Server) AskBot(ctx context.Context, rq *pb.AskBotRequest) (*pb.AskBotResponse, error) {
	rs, err := s.domain.AskBot(&domain.AskBotRequest{
		From: rq.From,
		Message:   rq.Message,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AskBotResponse{
		Found:  rs.Found,
		Answer: rs.Answer,
	}, nil
}

func (s *Server) SendPostFromBot(ctx context.Context, rq *pb.SendPostFromBotRequest) (*pb.SendPostFromBotResponse, error) {
	err := s.domain.SendPostFromBot(&domain.SendPostRequest{
		Message:   rq.Message,
		UserId:    rq.UserId,
		ChannelId: rq.ChannelId,
		Ephemeral: rq.Ephemeral,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SendPostFromBotResponse{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, rq *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.domain.DeleteUser(rq.MMUserId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{}, nil
}
