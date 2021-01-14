package domain

import "gitlab.medzdrav.ru/prototype/services/repository/storage"

func (s *serviceImpl) balanceFromDto(userId string, dtos []storage.Balance) *UserBalance {

	rs := &UserBalance{
		UserId:  userId,
		Balance: map[ServiceType]Balance{},
	}

	types := s.GetTypes()

	for _, d := range dtos {

		rs.Balance[types[d.ServiceTypeId]] = Balance{
			Available: d.Total - d.Delivered - d.Locked,
			Locked:    d.Locked,
			Total:     d.Total,
			Delivered: d.Delivered,
		}

	}

	return rs

}
