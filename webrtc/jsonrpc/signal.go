package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v3"
	"github.com/sourcegraph/jsonrpc2"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/meta"
	"sync"
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

type signal struct {
	peer             domain.Peer
	webrtc           domain.WebrtcService
	sync.Mutex
}

func newSignal(peer domain.Peer, webrtc domain.WebrtcService) *signal {
	return &signal{
		peer:     peer,
		webrtc:   webrtc,
		Mutex:    sync.Mutex{},
	}
}

func (s *signal) join(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (*webrtc.SessionDescription, error) {

	l := log.L().C(ctx).Pr("jsonrpc").Cmp("webrtc").Mth("join")

	var join Join
	err := json.Unmarshal(*req.Params, &join)
	if err != nil {
		l.E(err).Err("parsing offer")
		return nil, err
	}

	l.F(log.FF{"room": join.RoomId}).Dbg().TrcF("offer=%v", join.Offer)

	roomMeta, err := s.webrtc.GetOrCreateRoom(ctx, join.RoomId)
	if err != nil {
		l.E(err).Err("parsing offer")
		return nil, err
	}

	// Redirect to another node
	if meta.NodeId != roomMeta.NodeId {
		payload, _ := json.Marshal(roomMeta)
		// room exists on other node, let client know
		l.Warn("redirect")
		_ = conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    302,
			Message: string(payload),
		})
		return nil, nil
	}

	l.TrcF("offer=%v", join.Offer)
	answer, err := s.peer.Join(ctx, join.RoomId, &join.Offer)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}
	l.TrcF("answer=%s", answer.SDP)

	s.peer.OnOffer(func(offer *webrtc.SessionDescription) {
		l := log.L().C(ctx).Pr("jsonrpc").Cmp("webrtc").Mth("on-offer").F(log.FF{"room": join.RoomId}).Dbg().TrcF("%v", *offer)
		if err := conn.Notify(ctx, "offer", offer); err != nil {
			l.E(err).Err("sending offer")
		}
	})

	s.peer.OnIceCandidate(func(candidate *webrtc.ICECandidateInit, target int) {
		l := log.L().C(ctx).Pr("jsonrpc").Cmp("webrtc").Mth("on-ice-candidate").F(log.FF{"room": join.RoomId}).Dbg().TrcF("%v", candidate.Candidate)
		if err := conn.Notify(ctx, "trickle", Trickle{
			Candidate: *candidate,
			Target:    target,
		}); err != nil {
			l.E(err).Err("ice candidate")
		}
	})

	s.peer.OnICEConnectionStateChange(func(ss webrtc.ICEConnectionState) {
		l := log.L().C(ctx).Pr("jsonrpc").Cmp("webrtc").Mth("on-ice-conn-state").F(log.FF{"room": join.RoomId}).Inf(ss.String())
		if ss == webrtc.ICEConnectionStateFailed || ss == webrtc.ICEConnectionStateClosed {
			l.Inf("peer ice failed/closed, closing peer and websocket")
			s.peer.Close(ctx)
			conn.Close()
		}
	})

	return answer, nil
}

func (s *signal) offer(ctx context.Context, req *jsonrpc2.Request) (*webrtc.SessionDescription, error) {

	l := log.L().C(ctx).Pr("jsonrpc").Cmp("webrtc").Mth("offer").F(log.FF{"id": req.ID}).Dbg()

	var negotiation Negotiation
	err := json.Unmarshal(*req.Params, &negotiation)
	if err != nil {
		l.E(err).Err("parsing offer")
		return nil, err
	}

	answer, err := s.peer.Offer(&negotiation.Desc)
	if err != nil {
		l.E(err).Err("answer")
		return nil, err
	}
	l.TrcF("answer=%s", answer.SDP)

	return answer, nil
}

func (s *signal) answer(ctx context.Context, req *jsonrpc2.Request) error {

	l := log.L().C(ctx).Pr("jsonrpc").Cmp("webrtc").Mth("answer").F(log.FF{"id": req.ID}).Dbg()

	var negotiation Negotiation
	err := json.Unmarshal(*req.Params, &negotiation)
	if err != nil {
		l.E(err).Err("parsing offer")
		return err
	}

	err = s.peer.Answer(&negotiation.Desc)
	if err != nil {
		l.E(err).Err("set remote sdp")
		return err
	}

	return nil
}

func (s *signal) trickle(ctx context.Context, req *jsonrpc2.Request) error {

	l := log.L().C(ctx).Pr("jsonrpc").Cmp("webrtc").Mth("trickle").F(log.FF{"id": req.ID}).Dbg()

	var trickle Trickle
	err := json.Unmarshal(*req.Params, &trickle)
	if err != nil {
		l.E(err).Err("parsing offer")
		return err
	}

	err = s.peer.Trickle(trickle.Candidate, trickle.Target)
	if err != nil {
		l.E(err).Err()
		return err
	}

	return nil
}

// Handle incoming RPC call events like join, answer, offer and trickle
func (s *signal) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {

	ctx = kitContext.NewRequestCtx().
		WithRequestId(req.ID.String()).
		Webrtc().
		WithUser(s.peer.GetUserId(), s.peer.GetUsername()).
		ToContext(ctx)

	log.L().Pr("jsonrpc").C(ctx).Cmp("webrtc").Mth("handle").F(log.FF{"id": req.ID, "method": req.Method}).Dbg().TrcF("%v", req)

	s.Lock()
	defer s.Unlock()

	replyError := func(err error) {
		_ = conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    500,
			Message: fmt.Sprintf("%s", err),
		})
	}

	switch req.Method {

	case "join":
		answer, err := s.join(ctx, conn, req)
		if err != nil {
			replyError(err)
			break
		}

		if answer != nil {
			_ = conn.Reply(ctx, req.ID, answer)
			break
		}

	case "offer":
		answer, err := s.offer(ctx, req)
		if err != nil {
			replyError(err)
			break
		}

		if answer != nil {
			_ = conn.Reply(ctx, req.ID, answer)
			break
		}

	case "answer":
		err := s.answer(ctx, req)
		if err != nil {
			replyError(err)
			break
		}

	case "trickle":
		err := s.trickle(ctx, req)
		if err != nil {
			replyError(err)
			break
		}

	case "ping":
		_ = conn.Reply(ctx, req.ID, "pong")
		break
	}
}
