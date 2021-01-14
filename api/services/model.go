package services

type Balance struct {
	Available int `json:"available"`
	Delivered int `json:"delivered"`
	Locked    int `json:"locked"`
	Total     int `json:"total"`
}

type UserBalance struct {
	UserId  string `json:"userId"`
	Balance map[string]Balance `json:"balance"`
}

type ModifyUserBalanceRequest struct {
	ServiceTypeId string `json:"serviceTypeId"`
	Quantity      int `json:"quantity"`
}