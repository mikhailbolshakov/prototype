package impl

import (
	"context"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/webrtc/v3"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
)

type peerImpl struct {
	userId                       string
	username                     string
	sfuPeer                      *sfu.Peer
	onOfferEv                    domain.OnOfferEvent
	onIceCandidateEv             domain.OnIceCandidateEvent
	onIceConnectionStateChangeEv domain.OnICEConnectionStateChangeEvent
}

func newPeer(ctx context.Context, sessionProvider sfu.SessionProvider, userId, username string) domain.Peer {

	p := &peerImpl{
		userId: userId,
		username: username,
	}

	sfuPeer := sfu.NewPeer(sessionProvider)
	sfuPeer.OnOffer = p.onOffer
	sfuPeer.OnIceCandidate = p.onIceCandidate
	sfuPeer.OnICEConnectionStateChange = p.onICEConnectionStateChange
	p.sfuPeer = sfuPeer

	log.L().C(ctx).Cmp("webrtc").Mth("new-peer").Dbg("ok")

	return p

}

func (p *peerImpl) onOffer(o *webrtc.SessionDescription) {
	log.L().Cmp("webrtc").Mth("peer-onoffer").Dbg().TrcF("%v", *o)
	if p.onOfferEv != nil {
		p.onOfferEv(o)
	}
}

func (p *peerImpl) onIceCandidate(c *webrtc.ICECandidateInit, t int) {
	log.L().Cmp("webrtc").Mth("peer-onice").Dbg().TrcF("%v", *c)
	if p.onIceCandidateEv != nil {
		p.onIceCandidateEv(c, t)
	}
}

func (p *peerImpl) onICEConnectionStateChange(s webrtc.ICEConnectionState) {
	log.L().Cmp("webrtc").Mth("on-ice-state").DbgF("%v", s)
	if p.onIceConnectionStateChangeEv != nil {
		p.onIceConnectionStateChangeEv(s)
	}
}

func (p *peerImpl) Join(ctx context.Context, roomId string, offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {

	l := log.L().Cmp("webrtc").Mth("peer-join").F(log.FF{"room": roomId, "usr": p.userId}).Dbg().TrcF("%v", *offer)

	err := p.sfuPeer.Join(roomId, p.userId)
	if err != nil {
		return nil, err
	}

	answer, err := p.sfuPeer.Answer(*offer)
	if err != nil {
		return nil, err
	}
	l.F(log.FF{"answer": answer}).Trc()

	return answer, nil
}

func (p *peerImpl) Offer(sdp *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	log.L().Cmp("webrtc").Mth("peer-offer").Dbg().TrcF("%v", *sdp)
	return p.sfuPeer.Answer(*sdp)
}

func (p *peerImpl) Answer(sdp *webrtc.SessionDescription) error {
	log.L().Cmp("webrtc").Mth("peer-answer").Dbg().TrcF("%v", *sdp)
	return p.sfuPeer.SetRemoteDescription(*sdp)
}

func (p *peerImpl) Trickle(candidate webrtc.ICECandidateInit, target int) error {
	log.L().Cmp("webrtc").Mth("peer-trickle").Dbg().TrcF("%v", candidate)
	return p.sfuPeer.Trickle(candidate, target)
}

func (p *peerImpl) Close(ctx context.Context) {
	log.L().Cmp("webrtc").Mth("peer-close").Dbg()
	p.sfuPeer.Close()
}

func (p *peerImpl) OnOffer(e domain.OnOfferEvent) {
	p.onOfferEv = e
}

func (p *peerImpl) OnIceCandidate(e domain.OnIceCandidateEvent) {
	p.onIceCandidateEv = e
}

func (p *peerImpl) OnICEConnectionStateChange(e domain.OnICEConnectionStateChangeEvent) {
	p.onIceConnectionStateChangeEv = e
}

func (p *peerImpl) GetUserId() string {
	return p.userId
}

func (p *peerImpl) GetUsername() string {
	return p.username
}
