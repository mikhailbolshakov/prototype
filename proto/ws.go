package proto

const (
	// messages to interact with WebRTC signaling
	WS_MESSAGE_TYPE_WEBRTC = "webrtc"
	// messages to interact with a chat
	WS_MESSAGE_TYPE_CHAT = "chat"
)

// WsMessage is a common template of websocket messages
// all messages must correspond the template to be understood by the socket hub and forwarded to a recipient
type WsMessage struct {
	Id            string      `json:"id"`     // Id (Optional) used to trace a request and correlate with a response
	CorrelationId string      `json:"corrId"` // CorrelationId (Optional) allows to correlate with another message. Normally, RequestId from the context
	MessageType   string      `json:"type"`   // MessageType type of message
	Data          interface{} `json:"data"`   // Data it has arbitrary format and should be understood by a recipient of the given message type
}
