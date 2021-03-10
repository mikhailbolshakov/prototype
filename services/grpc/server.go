package grpc

import (
	"context"
	"encoding/json"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"gitlab.medzdrav.ru/prototype/services/logger"
	"gitlab.medzdrav.ru/prototype/services/meta"
	"log"
)

type Server struct {
	host, port string
	*kitGrpc.Server
	balanceService domain.UserBalanceService
	deliveryService domain.DeliveryService
	pb.UnimplementedBalanceServiceServer
	pb.UnimplementedDeliveryServiceServer
}

func New(balanceService domain.UserBalanceService,
		 deliveryService domain.DeliveryService) *Server {

	s := &Server{
		balanceService: balanceService,
		deliveryService: deliveryService,
	}

	// grpc server
	gs, err := kitGrpc.NewServer(meta.ServiceCode, logger.LF())
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterBalanceServiceServer(s.Srv, s)
	pb.RegisterDeliveryServiceServer(s.Srv, s)

	return s
}

func  (s *Server) Init(c *config.Config) error {
	cfg := c.Services["services"]
	s.host = cfg.Grpc.Host
	s.port = cfg.Grpc.Port
	return nil
}


func (s *Server) ListenAsync() {

	go func () {
		err := s.Server.Listen(s.host, s.port)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Add(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	rs, err := s.balanceService.Add(ctx, &domain.ModifyBalanceRequest{
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      int(rq.Quantity),
	})
	if err != nil {
		return nil, err
	}

	return s.toUserBalancePb(rs), nil
}

func (s *Server) GetBalance(ctx context.Context, rq *pb.GetBalanceRequest) (*pb.UserBalance, error) {

	rs, err := s.balanceService.Get(ctx, &domain.GetBalanceRequest{
		UserId:        rq.UserId,
	})
	if err != nil {
		return nil, err
	}

	return s.toUserBalancePb(rs), nil
}

func (s *Server) WriteOff(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {

	rs, err := s.balanceService.WriteOff(ctx, &domain.ModifyBalanceRequest{
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      int(rq.Quantity),
	})
	if err != nil {
		return nil, err
	}

	return s.toUserBalancePb(rs), nil
}

func (s *Server) Lock(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	rs, err := s.balanceService.Lock(ctx, &domain.ModifyBalanceRequest{
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      int(rq.Quantity),
	})
	if err != nil {
		return nil, err
	}

	return s.toUserBalancePb(rs), nil
}

func (s *Server) CancelLock(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {

	rs, err := s.balanceService.Cancel(ctx, &domain.ModifyBalanceRequest{
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      int(rq.Quantity),
	})
	if err != nil {
		return nil, err
	}

	return s.toUserBalancePb(rs), nil
}

func (s *Server) GetDelivery(ctx context.Context, rq *pb.GetDeliveryRequest) (*pb.Delivery, error) {
	return s.toDeliveryPb(s.deliveryService.Get(ctx, rq.Id)), nil
}

func (s *Server) Cancel(ctx context.Context, rq *pb.CancelDeliveryRequest) (*pb.Delivery, error) {

	d, err := s.deliveryService.Cancel(ctx, rq.Id, kitGrpc.PbTSToTime(rq.CancelTime))
	if err != nil {
		return nil, err
	}
	return s.toDeliveryPb(d), nil
}

func (s *Server) Complete(ctx context.Context, rq *pb.CompleteDeliveryRequest) (*pb.Delivery, error) {
	d, err := s.deliveryService.Complete(ctx, rq.Id, kitGrpc.PbTSToTime(rq.CompleteTime))
	if err != nil {
		return nil, err
	}
	return s.toDeliveryPb(d), nil
}

func (s *Server) UpdateDetails(ctx context.Context, rq *pb.UpdateDetailsRequest) (*pb.Delivery, error) {

	var details map[string]interface{}
	err := json.Unmarshal(rq.Details, &details)
	if err != nil {
		return nil, err
	}

	d, err := s.deliveryService.UpdateDetails(ctx, rq.Id, details)
	if err != nil {
		return nil, err
	}
	return s.toDeliveryPb(d), nil
}

func (s *Server) Create(ctx context.Context, rq *pb.DeliveryRequest) (*pb.Delivery, error) {

	var details map[string]interface{}
	err := json.Unmarshal(rq.Details, &details)
	if err != nil {
		return nil, err
	}

	d, err := s.deliveryService.Delivery(ctx, &domain.DeliveryRequest{
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Details:       details,
	})
	if err != nil {
		return nil, err
	}

	return s.toDeliveryPb(d), nil
}