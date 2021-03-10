package domain

import (
	"context"
	"gitlab.medzdrav.ru/prototype/proto/config"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"time"
)

type Storage interface {
	CreateBalance(ctx context.Context, b *BalanceItem) (*BalanceItem, error)
	UpdateBalance(ctx context.Context, b *BalanceItem) (*BalanceItem, error)
	GetBalance(ctx context.Context, userId string, at *time.Time) ([]*BalanceItem, error)
	GetBalanceForServiceType(ctx context.Context, userId string, serviceTypeId string, at *time.Time) ([]*BalanceItem, error)
	GetTypes(ctx context.Context) []ServiceType
	CreateDelivery(ctx context.Context, d *Delivery) (*Delivery, error)
	UpdateDelivery(ctx context.Context, d *Delivery) (*Delivery, error)
	UpdateDetails(ctx context.Context, deliveryId string, details map[string]interface{}) (*Delivery, error)
	GetDelivery(ctx context.Context, id string) *Delivery
}

type UserService interface {
	Get(ctx context.Context, id string) *pb.User
}

type ConfigService interface {
	Get() (*config.Config, error)
}

type BpService interface {
	StartProcess(ctx context.Context, processId string, vars map[string]interface{}) (string, error)
}