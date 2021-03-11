package webrtc

import "time"

type CreateRoomRequest struct {
	ChannelId string `json:"channelId"`
}

type RoomParticipants struct {
	UserId   string     `json:"userId"`
	Username string     `json:"username"`
	JoinedAt time.Time  `json:"joinedAt"`
	LeaveAt  *time.Time `json:"leaveAt,omitempty"`
}

type RoomDetails struct {
	ChannelId    string              `json:"channelId"`
	Participants []*RoomParticipants `json:"participants"`
}

type Room struct {
	Id       string       `json:"id"`
	OpenedAt *time.Time   `json:"openedAt,omitempty"`
	ClosedAt *time.Time   `json:"closedAt,omitempty"`
	Details  *RoomDetails `json:"details"`
}
