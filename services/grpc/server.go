package grpc

import (
	"context"
	"encoding/json"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	"gitlab.medzdrav.ru/prototype/services/domain"
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
	gs, err := kitGrpc.NewServer()
	if err != nil {
		panic(err)
	}
	s.Server = gs
	pb.RegisterBalanceServiceServer(s.Srv, s)
	pb.RegisterDeliveryServiceServer(s.Srv, s)

	return s
}

func  (s *Server) Init(c *kitConfig.Config) error {
	usersCfg := c.Services["services"]
	s.host = usersCfg.Grpc.Host
	s.port = usersCfg.Grpc.Port
	return nil
}


func (s *Server) ListenAsync() {

	go func () {
		err := s.Server.Listen("localhost", "50054")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Add(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {
	rs, err := s.balanceService.Add(&domain.ModifyBalanceRequest{
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

	rs, err := s.balanceService.Get(&domain.GetBalanceRequest{
		UserId:        rq.UserId,
	})
	if err != nil {
		return nil, err
	}

	return s.toUserBalancePb(rs), nil
}

func (s *Server) WriteOff(ctx context.Context, rq *pb.ChangeServicesRequest) (*pb.UserBalance, error) {

	rs, err := s.balanceService.WriteOff(&domain.ModifyBalanceRequest{
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
	rs, err := s.balanceService.Lock(&domain.ModifyBalanceRequest{
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

	rs, err := s.balanceService.Cancel(&domain.ModifyBalanceRequest{
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
	return s.toDeliveryPb(s.deliveryService.Get(rq.Id)), nil
}

func (s *Server) Cancel(ctx context.Context, rq *pb.CancelDeliveryRequest) (*pb.Delivery, error) {

	d, err := s.deliveryService.Cancel(rq.Id, kitGrpc.PbTSToTime(rq.CancelTime))
	if err != nil {
		return nil, err
	}
	return s.toDeliveryPb(d), nil
}

func (s *Server) Complete(ctx context.Context, rq *pb.CompleteDeliveryRequest) (*pb.Delivery, error) {
	d, err := s.deliveryService.Complete(rq.Id, kitGrpc.PbTSToTime(rq.CompleteTime))
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

	d, err := s.deliveryService.UpdateDetails(rq.Id, details)
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

	d, err := s.deliveryService.Delivery(&domain.DeliveryRequest{
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Details:       details,
	})
	if err != nil {
		return nil, err
	}

	return s.toDeliveryPb(d), nil
}