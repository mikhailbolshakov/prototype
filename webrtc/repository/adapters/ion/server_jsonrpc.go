package ion

import (
	"context"
	"encoding/json"
	"fmt"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"sync"

	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/webrtc/v3"
	"github.com/sourcegraph/jsonrpc2"
)

// Join message sent when initializing a peer connection
type Join struct {
	RoomId string                    `json:"sid"`
	Offer  webrtc.SessionDescription `json:"offer"`
}

// Negotiation message sent when renegotiating the peer connection
type Negotiation struct {
	Desc webrtc.SessionDescription `json:"desc"`
}

// Trickle message sent when renegotiating the peer connection
type Trickle struct {
	Target    int                     `json:"target"`
	Candidate webrtc.ICECandidateInit `json:"candidate"`
}

type JSONSignal struct {
	mu sync.Mutex
	c  coordinator
	*sfu.Peer
}

func (p *JSONSignal) join(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (*webrtc.SessionDescription, error) {

	l := log.L().C(ctx).Pr("jsonrpc").Cmp("ion").Mth("join").F(log.FF{"id": req.ID}).Dbg()

	var join Join
	err := json.Unmarshal(*req.Params, &join)
	if err != nil {
		l.E(err).Err("parsing offer")
		return nil, err
	}

	meta, err := p.c.getOrCreateRoom(ctx, join.RoomId)
	if err != nil {
		l.E(err).Err("parsing offer")
		return nil, err
	}

	if meta.Redirect {
		payload, _ := json.Marshal(meta)
		// room exists on other node, let client know
		l.Warn("redirect")
		_ = conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    302,
			Message: string(payload),
		})
		return nil, nil
	}

	l.TrcF("offer=%v", join.Offer)
	answer, err := p.Join(join.RoomId, join.Offer)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}
	l.TrcF("answer=%s", answer.SDP)

	p.OnOffer = func(offer *webrtc.SessionDescription) {
		l := log.L().Pr("jsonrpc").Cmp("ion").Mth("on-offer").F(log.FF{"room": join.RoomId}).Dbg().TrcF("%v", offer)
		if err := conn.Notify(ctx, "offer", offer); err != nil {
			l.E(err).Err("sending offer")
		}
	}
	p.OnIceCandidate = func(candidate *webrtc.ICECandidateInit, target int) {
		l := log.L().Pr("jsonrpc").Cmp("ion").Mth("on-ice-candidate").F(log.FF{"room": join.RoomId}).Dbg().TrcF("%v", candidate.Candidate)
		if err := conn.Notify(ctx, "trickle", Trickle{
			Candidate: *candidate,
			Target:    target,
		}); err != nil {
			l.E(err).Err("ice candidate")
		}
	}
	p.OnICEConnectionStateChange = func(s webrtc.ICEConnectionState) {
		l := log.L().Pr("jsonrpc").Cmp("ion").Mth("on-ice-conn-state").F(log.FF{"room": join.RoomId}).Inf(s.String())
		if s == webrtc.ICEConnectionStateFailed || s == webrtc.ICEConnectionStateClosed {
			l.Inf("peer ice failed/closed, closing peer and websocket")
			p.Close()
			conn.Close()
		}
	}

	return answer, nil
}

func (p *JSONSignal) offer(ctx context.Context, req *jsonrpc2.Request) (*webrtc.SessionDescription, error) {

	l := log.L().C(ctx).Pr("jsonrpc").Cmp("ion").Mth("offer").F(log.FF{"id": req.ID}).Dbg()

	var negotiation Negotiation
	err := json.Unmarshal(*req.Params, &negotiation)
	if err != nil {
		l.E(err).Err("parsing offer")
		return nil, err
	}

	answer, err := p.Answer(negotiation.Desc)
	if err != nil {
		l.E(err).Err("answer")
		return nil, err
	}
	l.TrcF("answer=%s", answer.SDP)

	return answer, nil
}

func (p *JSONSignal) answer(ctx context.Context, req *jsonrpc2.Request) error {

	l := log.L().C(ctx).Pr("jsonrpc").Cmp("ion").Mth("answer").F(log.FF{"id": req.ID}).Dbg()

	var negotiation Negotiation
	err := json.Unmarshal(*req.Params, &negotiation)
	if err != nil {
		l.E(err).Err("parsing offer")
		return err
	}

	err = p.SetRemoteDescription(negotiation.Desc)
	if err != nil {
		l.E(err).Err("set remote sdp")
		return err
	}

	return nil
}

func (p *JSONSignal) trickle(ctx context.Context, req *jsonrpc2.Request) error {

	l := log.L().C(ctx).Pr("jsonrpc").Cmp("ion").Mth("trickle").F(log.FF{"id": req.ID}).Dbg()

	var trickle Trickle
	err := json.Unmarshal(*req.Params, &trickle)
	if err != nil {
		l.E(err).Err("parsing offer")
		return err
	}

	err = p.Trickle(trickle.Candidate, trickle.Target)
	if err != nil {
		l.E(err).Err()
		return err
	}

	return nil
}

// Handle incoming RPC call events like join, answer, offer and trickle
func (p *JSONSignal) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {

	log.L().Pr("jsonrpc").Cmp("ion").Mth("handle").F(log.FF{"id": req.ID, "method": req.Method}).Dbg().TrcF("%v", req)

	ctx = kitContext.NewRequestCtx().WithRequestId(req.ID.String()).Webrtc().ToContext(ctx)

	p.mu.Lock()
	defer p.mu.Unlock()

	replyError := func(err error) {
		_ = conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    500,
			Message: fmt.Sprintf("%s", err),
		})
	}

	switch req.Method {

	case "join":
		answer, err := p.join(ctx, conn, req)
		if err != nil {
			replyError(err)
			break
		}

		if answer != nil {
			_ = conn.Reply(ctx, req.ID, answer)
			break
		}

	case "offer":
		answer, err := p.offer(ctx, req)
		if err != nil {
			replyError(err)
			break
		}

		if answer != nil {
			_ = conn.Reply(ctx, req.ID, answer)
			break
		}

	case "answer":
		err := p.answer(ctx, req)
		if err != nil {
			replyError(err)
			break
		}

	case "trickle":
		err := p.answer(ctx, req)
		if err != nil {
			replyError(err)
			break
		}

	case "ping":
		_ = conn.Reply(ctx, req.ID, "pong")
		break
	}
}
