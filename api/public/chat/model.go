package chat

type PostRequest struct {
	ChannelId string `json:"channelId"`
	Message   string `json:"message"`
}

type EphemeralPostRequest struct {
	PostRequest
	ToUserId string `json:"toUserId"`
}
