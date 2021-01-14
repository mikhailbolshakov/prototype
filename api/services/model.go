package services

import "time"

type Balance struct {
	ServiceTypeId string `json:"serviceType"`
	Available     int    `json:"available"`
	Delivered     int    `json:"delivered"`
	Locked        int    `json:"locked"`
	Total         int    `json:"total"`
}

type UserBalance struct {
	UserId  string    `json:"userId"`
	Balance []Balance `json:"balance"`
}

type ModifyUserBalanceRequest struct {
	ServiceTypeId string `json:"serviceTypeId"`
	Quantity      int    `json:"quantity"`
}

type DeliveryRequest struct {
	ServiceTypeId string `json:"serviceTypeId"`
	Details       map[string]interface{} `json:"details"`
}

type Delivery struct {
	Id            string `json:"id"`
	UserId        string `json:"userId"`
	ServiceTypeId string `json:"serviceTypeId"`
	Status        string `json:"status"`
	StartTime     *time.Time `json:"startTime"`
	FinishTime    *time.Time `json:"finishTime"`
	Details       map[string]interface{} `json:"details"`
}
