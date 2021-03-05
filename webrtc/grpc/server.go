package grpc

import (
	"encoding/json"
	"fmt"
	sfuPb "github.com/pion/ion-sfu/cmd/signal/grpc/proto"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/webrtc/v3"
	"gitlab.medzdrav.ru/prototype/kit"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	kitGrpc "gitlab.medzdrav.ru/prototype/kit/grpc"
	log "gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/meta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"sync"
)

type Server struct {
	host, port string
	// TODO: mutex per stream
	sync.Mutex
	*kitGrpc.Server
	webrtc domain.WebrtcService
	sfuPb.UnimplementedSFUServer
}

func New(domain domain.WebrtcService) *Server {

	s := &Server{webrtc: domain}

	// grpc server
	gs, err := kitGrpc.NewServer(meta.ServiceCode)
	if err != nil {
		panic(err)
	}
	s.Server = gs

	return s
}

func  (s *Server) Init(c *kitConfig.Config) error {

	cfg := c.Services["webrtc"]
	s.host = cfg.Grpc.Host
	s.port = cfg.Grpc.Port
	sfuPb.RegisterSFUServer(s.Srv, s)
	//sfuPb.RegisterSFUServer(s.Srv, sfuServer.NewServer(s.webrtc.GetSFU()))

	return nil
}

func (s *Server) ListenAsync() {

	go func () {
		err := s.Server.Listen(s.host, s.port)
		if err != nil {
			log.L().Pr("grpc").Mth("listen").E(err).Err()
			return
		}
	}()
}

func (s *Server) streamSend(stream sfuPb.SFU_SignalServer, r *sfuPb.SignalReply) error {
	s.Lock()
	defer s.Unlock()
	return stream.Send(r)
}

func (s *Server) streamSendErr(stream sfuPb.SFU_SignalServer, err error) error {
	return s.streamSend(stream, &sfuPb.SignalReply{
		Payload: &sfuPb.SignalReply_Error{
			Error: err.Error(),
		},
	})
}

func (s *Server) Signal(stream sfuPb.SFU_SignalServer) error {

	ctx := stream.Context()
	rq, ok := kitContext.Request(ctx)
	if !ok {
		user := "webrtc-" + kit.NewId()
		rq = kitContext.NewRequestCtx().Webrtc().WithNewRequestId().WithUser(user, user)
		ctx = rq.ToContext(ctx)
	}

	peer := s.webrtc.NewPeer(ctx, rq.GetUserId(), rq.GetUsername())

	l := log.L().Pr("grpc").Cmp("webrtc").Mth("signal").C(ctx).Dbg()

	for {
		in, err := stream.Recv()

		if err != nil {
			peer.Close(ctx)
			if err == io.EOF {
				return nil
			}
			errStatus, _ := status.FromError(err)
			if errStatus.Code() == codes.Canceled {
				return nil
			}
			l.E(err).ErrF("signal error %v %v", errStatus.Message(), errStatus.Code())
			return err
		}

		switch payload := in.Payload.(type) {
		case *sfuPb.SignalRequest_Join:

			l.Dbg("signal->join")

			// Notify user of new ice candidate
			peer.OnIceCandidate(func(candidate *webrtc.ICECandidateInit, target int) {

				if candidate == nil {
					return
				}

				ll := l.Clone().F(log.FF{"event": "on-ice", "peer": peer.GetUsername()}).Dbg().TrcF("%v", *candidate)

				bytes, err := json.Marshal(candidate)
				if err != nil {
					ll.E(err).St().Err()
					return
				}
				err = s.streamSend(stream,&sfuPb.SignalReply{
					Payload: &sfuPb.SignalReply_Trickle{
						Trickle: &sfuPb.Trickle{
							Init:   string(bytes),
							Target: sfuPb.Trickle_Target(target),
						},
					},
				})
				if err != nil {
					l.ErrF("OnIceCandidate send error %v ", err)
				}
			})

			// Notify user of new offer
			peer.OnOffer(func(o *webrtc.SessionDescription) {

				ll := l.Clone().F(log.FF{"event": "on-offer", "peer": peer.GetUsername()}).Dbg().TrcF("%v", *o)

				marshalled, err := json.Marshal(o)
				if err != nil {
					err = s.streamSendErr(stream, fmt.Errorf("offer sdp marshal error: %w", err))
					if err != nil {
						ll.E(err).St().Err("grpc send")
					}
					return
				}
				err = s.streamSend(stream, &sfuPb.SignalReply{
					Payload: &sfuPb.SignalReply_Description{
						Description: marshalled,
					},
				})

				if err != nil {
					ll.E(err).St().Err("negotiation")
				}
			})

			peer.OnICEConnectionStateChange( func(c webrtc.ICEConnectionState) {
				ll := l.Clone().F(log.FF{"event": "on-ice-state", "peer": peer.GetUsername()}).Dbg().TrcF("%v", c)

				err = s.streamSend(stream, &sfuPb.SignalReply{
					Payload: &sfuPb.SignalReply_IceConnectionState{
						IceConnectionState: c.String(),
					},
				})

				if err != nil {
					ll.E(err).St().Err("oniceconnectionstatechange")
				}
			})

			offer, err := s.toSdpDomain(payload.Join.Description)
			if err != nil {
				err = s.streamSendErr(stream, fmt.Errorf("join sdp unmarshal error: %w", err))
				if err != nil {
					l.E(err).Err("grpc send")
					return status.Errorf(codes.Internal, err.Error())
				}
			}
			l.TrcF("offer=%v", *offer)

			// get or create room
			roomMeta, err := s.webrtc.GetOrCreateRoom(ctx, payload.Join.Sid)
			if err != nil {
				err = s.streamSendErr(stream, fmt.Errorf("get or create room: %w", err))
				l.E(err).Err("grpc send")
				continue
			}

			// TODO: redirect
			if roomMeta.NodeId != meta.NodeId {
				l.WarnF("redirecting to %s node", roomMeta.NodeId)
				return nil
			}

			answer, err := peer.Join(ctx, payload.Join.Sid, offer)
			l.Dbg("join").TrcF("answer=%v", *answer)

			if err != nil {
				switch err {
				case sfu.ErrTransportExists:
					fallthrough
				case sfu.ErrOfferIgnored:
					err = s.streamSendErr(stream, fmt.Errorf("join error: %w", err))
					if err != nil {
						l.ErrF("grpc send error %v ", err)
						return status.Errorf(codes.Internal, err.Error())
					}
				default:
					return status.Errorf(codes.Unknown, err.Error())
				}
			}

			marshalled, err := json.Marshal(answer)
			if err != nil {
				l.E(err).St().Err("marshal")
				return status.Errorf(codes.Internal, fmt.Sprintf("sdp marshal error: %v", err))
			}

			// send answer
			err = s.streamSend(stream, &sfuPb.SignalReply{
				Id: in.Id,
				Payload: &sfuPb.SignalReply_Join{
					Join: &sfuPb.JoinReply{
						Description: marshalled,
					},
				},
			})

			if err != nil {
				l.E(err).St().Err("sending join response")
				return status.Errorf(codes.Internal, "join error %s", err)
			}

		case *sfuPb.SignalRequest_Description:

			sdp, err := s.toSdpDomain(payload.Description)
			if err != nil {
				err = s.streamSendErr(stream, fmt.Errorf("negotiate sdp unmarshal error: %w", err))
				if err != nil {
					l.ErrF("grpc send error %v ", err)
					return status.Errorf(codes.Internal, err.Error())
				}
			}

			if sdp.Type == webrtc.SDPTypeOffer {
				l.DbgF("signal->offer").TrcF("sdp=%v", sdp)
				answer, err := peer.Offer(sdp)
				l.TrcF("answer=%v", answer)
				if err != nil {
					switch err {
					case sfu.ErrNoTransportEstablished:
						fallthrough
					case sfu.ErrOfferIgnored:
						err = s.streamSendErr(stream, fmt.Errorf("negotiate answer error: %w", err))
						if err != nil {
							l.E(err).St().Err("grpc send")
							return status.Errorf(codes.Internal, err.Error())
						}
						continue
					default:
						l.E(err).St().Err("negotiate")
						return status.Errorf(codes.Unknown, fmt.Sprintf("negotiate error: %v", err))
					}
				}

				marshalled, err := json.Marshal(answer)
				if err != nil {
					err = s.streamSendErr(stream, fmt.Errorf("sdp marshal error: %w", err))
					if err != nil {
						l.E(err).St().Err("grpc send")
						return status.Errorf(codes.Internal, err.Error())
					}
				}

				err = s.streamSend(stream, &sfuPb.SignalReply{
					Id: in.Id,
					Payload: &sfuPb.SignalReply_Description{
						Description: marshalled,
					},
				})

				if err != nil {
					l.E(err).St().Err("negotiate")
					return status.Errorf(codes.Internal, fmt.Sprintf("negotiate error: %v", err))
				}

			} else if sdp.Type == webrtc.SDPTypeAnswer {
				l.DbgF("signal->answer").TrcF("sdp=%v", sdp)
				err := peer.Answer(sdp)
				if err != nil {
					switch err {
					case sfu.ErrNoTransportEstablished:
						err = s.streamSendErr(stream, fmt.Errorf("set remote description error: %w", err))
						if err != nil {
							l.E(err).St().Err("grpc send")
							return status.Errorf(codes.Internal, err.Error())
						}
					default:
						l.E(err).St().Err("negotiate")
						return status.Errorf(codes.Unknown, err.Error())
					}
				}
			}

		case *sfuPb.SignalRequest_Trickle:
			l.DbgF("signal->trickle")
			var candidate webrtc.ICECandidateInit
			err := json.Unmarshal([]byte(payload.Trickle.Init), &candidate)
			if err != nil {
				l.E(err).St().Err("parsing ice candidate")
				err = s.streamSendErr(stream, fmt.Errorf("unmarshal ice candidate error:  %w", err))
				if err != nil {
					l.E(err).St().Err("grpc send")
					return status.Errorf(codes.Internal, err.Error())
				}
				continue
			}

			l.TrcF("candidate=%v", candidate)
			err = peer.Trickle(candidate, int(payload.Trickle.Target))
			if err != nil {
				switch err {
				case sfu.ErrNoTransportEstablished:
					l.E(err).St().Err("transport")
					err = s.streamSendErr(stream, fmt.Errorf("trickle error:  %w", err))
					if err != nil {
						l.E(err).St().Err("grpc send")
						return status.Errorf(codes.Internal, err.Error())
					}
				default:
					l.E(err).St().Err("negotiate")
					return status.Errorf(codes.Unknown, fmt.Sprintf("negotiate error: %v", err))
				}
			}

		}
	}
}
