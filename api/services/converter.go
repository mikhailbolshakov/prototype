package services

import pb "gitlab.medzdrav.ru/prototype/proto/services"

func (c *controller) balanceFromPb(p *pb.UserBalance) *UserBalance {

	rs := &UserBalance{
		UserId:  p.UserId,
		Balance: map[string]Balance{},
	}

	for k, v := range p.Balance {
		rs.Balance[k] = Balance{
			Available: int(v.Available),
			Delivered: int(v.Delivered),
			Locked:    int(v.Locked),
			Total:     int(v.Total),
		}
	}

	return rs
}