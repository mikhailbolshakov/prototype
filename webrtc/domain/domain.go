package domain

import (
	"context"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/webrtc/v3"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"time"
)

type RoomMeta struct {
	Id       string
	Endpoint string
	NodeId   string
}

type OnOfferEvent func(*webrtc.SessionDescription)
type OnIceCandidateEvent func(*webrtc.ICECandidateInit, int)
type OnICEConnectionStateChangeEvent func(webrtc.ICEConnectionState)
type OnError func(error)

type Peer interface {
	Join(ctx context.Context, roomId string, offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error)
	Offer(sdp *webrtc.SessionDescription) (*webrtc.SessionDescription, error)
	Answer(sdp *webrtc.SessionDescription) error
	Trickle(candidate webrtc.ICECandidateInit, target int) error

	Close(ctx context.Context)

	OnOffer(e OnOfferEvent)
	OnIceCandidate(e OnIceCandidateEvent)
	OnICEConnectionStateChange(e OnICEConnectionStateChangeEvent)

	GetUserId() string
	GetUsername() string
}

type RoomRecorder interface {
	Close()
}

// Recording creates a new peer and connects to SFU to relay all the tracks
type Recording interface {
	Init(ctx context.Context, cfg *config.Config, webrtc WebrtcService) error
	NewRoomRecorder(ctx context.Context, roomId string) (RoomRecorder, error)
}

type WebrtcService interface {
	Init(ctx context.Context, cfg *config.Config) error
	NewPeer(ctx context.Context, userId, username string) Peer
	GetSFU() *sfu.SFU
	GetOrCreateRoom(ctx context.Context, roomId string) (*RoomMeta, error)
}

type RoomParticipants struct {
	PeerId   string     `json:"peerId"`
	UserId   string     `json:"userId"`
	Username string     `json:"username"`
	JoinedAt time.Time  `json:"joinedAt"`
	LeaveAt  *time.Time `json:"leaveAt"`
}

type RoomDetails struct {
	ChannelId    string              `json:"channelId"`
	Participants []*RoomParticipants `json:"participants"`
}

type Room struct {
	Id         string       `json:"id"`
	OpenedAt   *time.Time   `json:"openedAt"`
	ClosedAt   *time.Time   `json:"closedAt"`
	Details    *RoomDetails `json:"details"`
	CreatedAt  time.Time    `json:"createdAt"`
	ModifiedAt time.Time    `json:"modifiedAt"`
	DeletedAt  *time.Time   `json:"deletedAt,omitempty"`
}

type RoomService interface {
	Create(ctx context.Context, channelId string) (*Room, error)
	Get(ctx context.Context, roomId string) *Room
	Join(ctx context.Context, roomId, userId, username, peerId string) (*Room, error)
	Leave(ctx context.Context, roomId, peerId string) (*Room, error)
	Close(ctx context.Context, roomId string) (*Room, error)
}
