package grpc

import (
	"encoding/json"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	"gitlab.medzdrav.ru/prototype/services/domain"
)

func (s *Server) toUserBalancePb(d *domain.UserBalance) *pb.UserBalance {
	rs := &pb.UserBalance{
		UserId:  d.UserId,
		Balance: make(map[string]*pb.Balance),
	}

	for key, val := range d.Balance {

		rs.Balance[key.Id] = &pb.Balance{
			Available: int32(val.Available),
			Delivered: int32(val.Delivered),
			Locked:    int32(val.Locked),
			Total:     int32(val.Total),
		}
	}
	return rs
}

func (s *Server) toDeliveryPb(d *domain.Delivery) *pb.Delivery {

	if d == nil {
		return nil
	}

	detailsB, _ := json.Marshal(d.Details)

	return &pb.Delivery{
		Id:            d.Id,
		UserId:        d.UserId,
		ServiceTypeId: d.ServiceTypeId,
		Status:        d.Status,
		StartTime:     kitGrpc.TimeToPbTS(&d.StartTime),
		FinishTime:    kitGrpc.TimeToPbTS(d.FinishTime),
		Details:       detailsB,
	}
}
