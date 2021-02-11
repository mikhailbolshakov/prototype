package domain

import (
	"context"
	"time"
)

const (
	ST_EXPERT_ONLINE_CONSULTATION = "expert-online-consultation"
)

type ModifyBalanceRequest struct {
	UserId        string
	ServiceTypeId string
	Quantity      int
}

type ServiceType struct {
	Id           string
	Description  string
	DeliveryWfId string
}

type Balance struct {
	Available int
	Delivered int
	Locked    int
	Total     int
}

type BalanceItem struct {
	Id            string
	UserId        string
	ServiceTypeId string
	Total         int
	Delivered     int
	Locked        int
}

type UserBalance struct {
	UserId  string
	Balance map[ServiceType]Balance
}

type GetBalanceRequest struct {
	UserId string
}

type DeliveryRequest struct {
	UserId        string
	ServiceTypeId string
	Details map[string]interface{}
}

type Delivery struct {
	Id string
	UserId string
	ServiceTypeId string
	Status string
	StartTime time.Time
	FinishTime *time.Time
	Details map[string]interface{}
}

type UserBalanceService interface {
	// get available service types
	GetTypes(ctx context.Context) map[string]ServiceType
	// adds service to balance
	Add(ctx context.Context, rq *ModifyBalanceRequest) (*UserBalance, error)
	// requests a balance
	Get(ctx context.Context, rq *GetBalanceRequest) (*UserBalance, error)
	// write off services
	WriteOff(ctx context.Context, rq *ModifyBalanceRequest) (*UserBalance, error)
	// lock service
	Lock(ctx context.Context, rq *ModifyBalanceRequest) (*UserBalance, error)
	// cancel locked service
	Cancel(ctx context.Context, rq *ModifyBalanceRequest) (*UserBalance, error)
}

type DeliveryService interface {
	// delivery user service
	Delivery(ctx context.Context, rq *DeliveryRequest) (*Delivery, error)
	Complete(ctx context.Context, deliveryId string, finishTime *time.Time) (*Delivery, error)
	Cancel(ctx context.Context, deliveryId string, cancelTime *time.Time) (*Delivery, error)
	Get(ctx context.Context, deliveryId string) *Delivery
	UpdateDetails(ctx context.Context, deliveryId string, details map[string]interface{}) (*Delivery, error)
}
