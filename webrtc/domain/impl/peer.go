package impl

import (
	"context"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/webrtc/v3"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/logger"
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

	p.l().C(ctx).Mth("new-peer").Dbg("ok")

	return p

}

func (p *peerImpl) l() log.CLogger {
	return logger.L().Cmp("webrtc-peer")
}

func (p *peerImpl) onOffer(o *webrtc.SessionDescription) {
	p.l().Mth("peer-onoffer").Dbg().TrcF("%v", *o)
	if p.onOfferEv != nil {
		p.onOfferEv(o)
	}
}

func (p *peerImpl) onIceCandidate(c *webrtc.ICECandidateInit, t int) {
	p.l().Mth("peer-onice").Dbg().TrcF("%v", *c)
	if p.onIceCandidateEv != nil {
		p.onIceCandidateEv(c, t)
	}
}

func (p *peerImpl) onICEConnectionStateChange(s webrtc.ICEConnectionState) {
	p.l().Mth("on-ice-state").DbgF("%v", s)
	if p.onIceConnectionStateChangeEv != nil {
		p.onIceConnectionStateChangeEv(s)
	}
}

func (p *peerImpl) Join(ctx context.Context, roomId string, offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {

	l := p.l().Mth("peer-join").F(log.FF{"room": roomId, "usr": p.userId}).Dbg().TrcF("%v", *offer)

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
	p.l().Mth("peer-offer").Dbg().TrcF("%v", *sdp)
	return p.sfuPeer.Answer(*sdp)
}

func (p *peerImpl) Answer(sdp *webrtc.SessionDescription) error {
	p.l().Mth("peer-answer").Dbg().TrcF("%v", *sdp)
	return p.sfuPeer.SetRemoteDescription(*sdp)
}

func (p *peerImpl) Trickle(candidate webrtc.ICECandidateInit, target int) error {
	p.l().Mth("peer-trickle").Dbg().TrcF("%v", candidate)
	return p.sfuPeer.Trickle(candidate, target)
}

func (p *peerImpl) Close(ctx context.Context) {
	p.l().Mth("peer-close").Dbg()
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
