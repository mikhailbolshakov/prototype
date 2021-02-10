package domain

import (
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
	GetTypes() map[string]ServiceType
	// adds service to balance
	Add(rq *ModifyBalanceRequest) (*UserBalance, error)
	// requests a balance
	Get(rq *GetBalanceRequest) (*UserBalance, error)
	// write off services
	WriteOff(rq *ModifyBalanceRequest) (*UserBalance, error)
	// lock service
	Lock(rq *ModifyBalanceRequest) (*UserBalance, error)
	// cancel locked service
	Cancel(rq *ModifyBalanceRequest) (*UserBalance, error)
}

type DeliveryService interface {
	// delivery user service
	Delivery(rq *DeliveryRequest) (*Delivery, error)
	Complete(deliveryId string, finishTime *time.Time) (*Delivery, error)
	Cancel(deliveryId string, cancelTime *time.Time) (*Delivery, error)
	Get(deliveryId string) *Delivery
	UpdateDetails(deliveryId string, details map[string]interface{}) (*Delivery, error)
}
