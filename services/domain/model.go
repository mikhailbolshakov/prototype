package domain

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
