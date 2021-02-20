package webrtc

type WsSignalRequest struct {
	Method  string                 `json:"methods"`
	Payload map[string]interface{} `json:"payload"`
}

type WsRoomSignal struct {

}