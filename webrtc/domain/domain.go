package domain

import (
	"context"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/webrtc/v3"
	"gitlab.medzdrav.ru/prototype/proto/config"
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
