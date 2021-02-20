package domain

const (

	SIGNAL_METHOD_FROM_CLIENT_JOIN    = "join-rq"
	SIGNAL_METHOD_FROM_CLIENT_OFFER   = "offer"
	SIGNAL_METHOD_FROM_CLIENT_ANSWER  = "answer"
	SIGNAL_METHOD_FROM_CLIENT_LEAVE   = "leave"
	SIGNAL_METHOD_FROM_CLIENT_TRICKLE = "trickle"

	SIGNAL_METHOD_TO_CLIENT_ROOM_RS = "room-rs"
	SIGNAL_METHOD_TO_CLIENT_JOIN_RS = "join-rs"
	SIGNAL_METHOD_TO_CLIENT_OFFER   = "offer"
	SIGNAL_METHOD_TO_CLIENT_ANSWER  = "answer"
	SIGNAL_METHOD_TO_CLIENT_TRICKLE = "trickle"
)

type WebRtcSignal struct {
	Method  string                 `json:"method"`
	Payload map[string]interface{} `json:"payload"`
}

type WebRtsSignalCreateRoomRequest struct {
	UserId    string `json:"userId"`    // UserId - user who initiated this room creation
	ChannelId string `json:"channelId"` // ChannelId (Optional) - populated if the room is created from a chat
	SDPOffer  string `json:"sdp"`       // SDPOffer - webrtc Offer session descriptor obtained on client side
}

type WebRtcSignalCreateRoomResponse struct {
	RoomId    string `json:"roomId"`
	SDPAnswer string `json:"sdp"` // SDPAnswer - webrtc Answer session descriptor
}

type WebRtsSignalOffer struct {
	UserId    string `json:"userId"`    // UserId - user who initiated this room creation
	ChannelId string `json:"channelId"` // ChannelId (Optional) - populated if the room is created from a chat
	SDPOffer  string `json:"sdp"`       // SDPOffer - webrtc Offer session descriptor obtained on client side
}

type WebRtcSignalAnswer struct {
	RoomId    string `json:"roomId"`
	SDPAnswer string `json:"sdp"` // SDPAnswer - webrtc Answer session descriptor
}

type JoinRoomRequest struct {
	RoomId   string
	UserId   string
	SDPOffer string
}

type JoinRoomResponse struct {
	SDPAnswer string // SDPAnswer - webrtc Answer session descriptor
}

// RoomHub is responsible for managing webrtc rooms
type RoomHub interface {
	WebrtcWsMessageHandler(msg []byte) error
}
