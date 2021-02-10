package domain

import (
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"time"
)

type Storage interface {
	CreateBalance(b *BalanceItem) (*BalanceItem, error)
	UpdateBalance(b *BalanceItem) (*BalanceItem, error)
	GetBalance(userId string, at *time.Time) ([]*BalanceItem, error)
	GetBalanceForServiceType(userId string, serviceTypeId string, at *time.Time) ([]*BalanceItem, error)
	GetTypes() []ServiceType
	CreateDelivery(d *Delivery) (*Delivery, error)
	UpdateDelivery(d *Delivery) (*Delivery, error)
	UpdateDetails(deliveryId string, details map[string]interface{}) (*Delivery, error)
	GetDelivery(id string) *Delivery
}

type UserService interface {
	Get(id string) *pb.User
}

type ConfigService interface {
	Get() (*kit.Config, error)
}

type BpService interface {
	StartProcess(processId string, vars map[string]interface{}) (string, error)
}