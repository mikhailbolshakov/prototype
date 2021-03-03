package domain

import (
	"context"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/webrtc/v3"
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
	Join(ctx context.Context, roomId, userId string, offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error)
	Offer(sdp *webrtc.SessionDescription) (*webrtc.SessionDescription, error)
	Answer(sdp *webrtc.SessionDescription) error
	Trickle(candidate webrtc.ICECandidateInit, target int) error

	Close(ctx context.Context)

	SetOnOffer(e OnOfferEvent)
	SetOnIceCandidate(e OnIceCandidateEvent)
	SetOnICEConnectionStateChange(e OnICEConnectionStateChangeEvent)
}

type WebrtcService interface {
	NewPeer(ctx context.Context) Peer
	GetSFU() *sfu.SFU
	GetOrCreateRoom(ctx context.Context, roomId string) (*RoomMeta, error)
}
