package proto

const (
	// messages sent to this topic are forwarded to websocket connections recognized by UserId
	QUEUE_TOPIC_OUTGOING_WS_EVENT    = "ws.event.outgoing"
	// is a template topic that turns into a real topic by concatenation with websocket message type
	// e.g. "ws.event.incoming.chat"
	QUEUE_TOPIC_INCOMING_WS_TEMPLATE = "ws.event.incoming.%s"
)

// OutgoingWsEventQueueMessagePayload defines a queue message format to be sent to QUEUE_TOPIC_OUTGOING_WS_EVENT topic
type OutgoingWsEventQueueMessagePayload struct {
	UserId  string     `json:"userId"`    // UserId - to whom this message is supposed to be sent
	WsEvent *WsMessage `json:"wsMessage"` // WsEvent - message Payload to be sent to websocket
}
