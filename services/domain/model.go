package domain

import "time"

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
