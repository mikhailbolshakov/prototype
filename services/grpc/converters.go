package grpc

import (
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	"gitlab.medzdrav.ru/prototype/services/domain"
)

func userBalanceFromDomain(d *domain.UserBalance) *pb.UserBalance {
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
