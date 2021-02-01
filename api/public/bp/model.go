package bp

type StartProcessRequest struct {
	ProcessId string                 `json:"processId"`
	Vars      map[string]interface{} `json:"vars"`
}

type StartProcessResponse struct {
	Id string `json:"id"`
}
