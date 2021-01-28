package services

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
)

func (c *ctrlImpl) balanceFromPb(p *pb.UserBalance) *UserBalance {

	rs := &UserBalance{
		UserId:  p.UserId,
		Balance: []Balance{},
	}

	for k, v := range p.Balance {
		rs.Balance = append(rs.Balance, Balance{
			ServiceTypeId: k,
			Available: int(v.Available),
			Delivered: int(v.Delivered),
			Locked:    int(v.Locked),
			Total:     int(v.Total),
		})
	}

	return rs
}

func (c *ctrlImpl) deliveryFromPb(p *pb.Delivery) *Delivery {

	var details map[string]interface{}
	err := json.Unmarshal(p.Details, &details)
	if err != nil {
		return nil
	}

	rs := &Delivery{
		Id:            p.Id,
		UserId:        p.UserId,
		ServiceTypeId: p.ServiceTypeId,
		Status:        p.Status,
		StartTime:     grpc.PbTSToTime(p.StartTime),
		FinishTime:    grpc.PbTSToTime(p.FinishTime),
		Details:       details,
	}

	return rs
}