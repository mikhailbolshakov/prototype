package impl

import (
	"context"
	"encoding/json"
	"fmt"
	avp "github.com/pion/ion-avp/pkg"
	"github.com/pion/ion-avp/pkg/elements"
	sfu "github.com/pion/ion-sfu/cmd/signal/grpc/proto"
	"github.com/pion/webrtc/v3"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	avpCstm "gitlab.medzdrav.ru/prototype/webrtc/domain/impl/avp"
	"gitlab.medzdrav.ru/prototype/webrtc/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"path"
)

const RECORDING_PEER_CODE = "avp"

type recordingIml struct {
	cfg         *config.Config
	sfuClient   sfu.SFUClient
	webrtc      domain.WebrtcService
}

type webmSaverRoomRecorderImpl struct {
	cancelTransportFn context.CancelFunc
}

func NewRecording() domain.Recording {
	return &recordingIml{}
}

func (r *recordingIml) l() log.CLogger {
	return logger.L().Cmp("webrtc-rec")
}

func (r *recordingIml) Init(ctx context.Context, cfg *config.Config, webrtc domain.WebrtcService) error {

	r.cfg = cfg
	r.webrtc = webrtc

	if cfg.Webrtc.Recording.File.Enabled {
		if err := os.MkdirAll(cfg.Webrtc.Recording.File.Path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// createWebmSaver - creates pipeline filters (webm + filesaver)
func (r *recordingIml) createWebmSaver() avp.ElementFun {

	return func(sid, pid, tid string, config []byte) avp.Element {
		filewriter := elements.NewFileWriter(
			path.Join(r.cfg.Webrtc.Recording.File.Path, fmt.Sprintf("%s-%s.webm", sid, pid)),
			4096,
		)
		webm := elements.NewWebmSaver()
		webm.Attach(filewriter)
		return webm
	}

}

//	NewRoomRecorderGRPC - connects to SFU via standard gRPC protocol (not used now because direct connection also works, but gRPC approach looks more reliable)
//	don't forget initialize gRPC connection in Init method
//
//	client, err := kitGrpc.NewClient(cfg.Services["webrtc"].Grpc.Host, cfg.Services["webrtc"].Grpc.Port)
//	if err != nil {
//		return err
//	}
//	r.sfuClient = sfu.NewSFUClient(client.Conn)
func (r *recordingIml) NewRoomRecorderGRPC(ctx context.Context, roomId string) (domain.RoomRecorder, error) {

	ctx, cancel := context.WithCancel(ctx)

	l := r.l().C(ctx).Mth("new-room-recorder").F(log.FF{"room": roomId}).Dbg()

	webrtcTransport := avpCstm.NewAvpPeer(roomId, r.cfg.Webrtc.Avp, r.createWebmSaver())
	webrtcTransport.OnClose(func() {
		cancel()
	})

	sfuStream, err := r.sfuClient.Signal(ctx)
	if err != nil {
		return nil, err
	}

	offer, err := webrtcTransport.CreateOffer()
	if err != nil {
		//log.Errorf("Error creating offer: %v", err)
		return nil, err
	}

	marshalled, err := json.Marshal(offer)
	if err != nil {
		return nil, err
	}

	l.TrcF("offer=%v", offer)

	err = sfuStream.Send(
		&sfu.SignalRequest{
			Payload: &sfu.SignalRequest_Join{
				Join: &sfu.JoinRequest{
					Sid:         roomId,
					Description: marshalled,
				},
			},
		},
	)

	if err != nil {
		//log.Errorf("Error sending publish request: %v", err)
		return nil, err
	}

	webrtcTransport.OnICECandidate(func(c *webrtc.ICECandidate, target int) {

		ll := l.Clone().F(log.FF{"event": "on-ice-candidate"})

		if c == nil {
			ll.Dbg("done")
			// Gathering done
			return
		}
		ll.TrcF("%v", *c)
		bytes, err := json.Marshal(c.ToJSON())
		if err != nil {
			ll.E(err).St().Err()
			return
		}
		err = sfuStream.Send(&sfu.SignalRequest{
			Payload: &sfu.SignalRequest_Trickle{
				Trickle: &sfu.Trickle{
					Init:   string(bytes),
					Target: sfu.Trickle_Target(target),
				},
			},
		})
		if err != nil {
			ll.E(err).St().Err()
		}
	})

	go func() {
		// Handle sfu stream messages

		for {
			ll := l.Clone().F(log.FF{"event": "on-stream"})

			res, err := sfuStream.Recv()

			if err != nil {
				ll.E(err).St().Err()
				if err == io.EOF {
					// WebRTC Transport closed
					ll.Inf("closed")
					err = sfuStream.CloseSend()
					if err != nil {
						ll.E(err).St().Err()
					}
					return
				}

				errStatus, _ := status.FromError(err)
				if errStatus.Code() == codes.Canceled {
					err = sfuStream.CloseSend()
					if err != nil {
						ll.E(err).St().Err()
					}
					return
				}
				return
			}

			switch payload := res.Payload.(type) {
			case *sfu.SignalReply_Join:
				// Set the remote SessionDescription
				ll.DbgF("answer").TrcF("%v", string(payload.Join.Description))

				var sdp webrtc.SessionDescription
				err := json.Unmarshal(payload.Join.Description, &sdp)
				if err != nil {
					ll.E(err).St().Err("unmarshal")
					return
				}

				if err = webrtcTransport.SetRemoteDescription(sdp); err != nil {
					ll.E(err).St().Err("join")
					return
				}

			case *sfu.SignalReply_Description:

				var sdp webrtc.SessionDescription
				err := json.Unmarshal(payload.Description, &sdp)
				if err != nil {
					ll.E(err).St().Err("unmarshal")
					return
				}

				if sdp.Type == webrtc.SDPTypeOffer {

					ll.Dbg("offer").TrcF("%v", sdp)

					var answer webrtc.SessionDescription
					answer, err = webrtcTransport.Answer(sdp)
					if err != nil {
						ll.E(err).St().Err("negotiation")
						continue
					}

					marshalled, err = json.Marshal(answer)
					if err != nil {
						ll.E(err).St().Err("negotiation")
						continue
					}

					err = sfuStream.Send(&sfu.SignalRequest{
						Payload: &sfu.SignalRequest_Description{
							Description: marshalled,
						},
					})

					if err != nil {
						ll.E(err).St().Err("negotiation")
						continue
					}

				} else if sdp.Type == webrtc.SDPTypeAnswer {

					ll.Dbg("answer").TrcF("%v", sdp)
					err = webrtcTransport.SetRemoteDescription(sdp)

					if err != nil {
						ll.E(err).St().Err("negotiation")
						continue
					}
				}
			case *sfu.SignalReply_Trickle:

				ll.Dbg("candidate").TrcF("%v", payload.Trickle.Init)
				var candidate webrtc.ICECandidateInit
				_ = json.Unmarshal([]byte(payload.Trickle.Init), &candidate)
				err := webrtcTransport.AddICECandidate(candidate, int(payload.Trickle.Target))
				if err != nil {
					ll.E(err).St().Err("trickle")
				}
			}
		}
	}()

	return &webmSaverRoomRecorderImpl{
		cancelTransportFn: cancel,
	}, nil
}

func (r *recordingIml) NewRoomRecorder(ctx context.Context, roomId string) (domain.RoomRecorder, error) {

	l := r.l().C(ctx).Mth("new-room-recorder").F(log.FF{"room": roomId})

	ctx, cancel := context.WithCancel(ctx)

	// create sfu-peer
	peer := r.webrtc.NewPeer(ctx, RECORDING_PEER_CODE, RECORDING_PEER_CODE)

	// create avp-peer
	avpPeer := avpCstm.NewAvpPeer(roomId, r.cfg.Webrtc.Avp, r.createWebmSaver())

	avpPeer.OnClose(func() {
		l.Clone().F(log.FF{"peer": "avp", "event": "on-close"}).Dbg("on close")
		peer.Close(ctx)
		cancel()
	})

	peer.OnICEConnectionStateChange(func(c webrtc.ICEConnectionState) {
		l.Clone().F(log.FF{"peer": peer.GetUsername(), "event": "on-ice-state"}).Dbg(c)
		avpPeer.OnIceConnectionStateChanged(c)
	})

	offer, err := avpPeer.CreateOffer()
	if err != nil {
		return nil, err
	}

	peer.OnOffer(func(sdp *webrtc.SessionDescription) {

		ll := l.Clone().F(log.FF{"peer": peer.GetUsername(), "event": "on-offer"}).Dbg().TrcF("sdp=%v", *sdp)

		answer, err := avpPeer.Answer(*sdp)
		if err != nil {
			ll.E(err).St().Err()
			return
		}

		ll.DbgF("avp answer = %v", answer)

		// to avoid lock another goroutine is used
		go func() {
			err := peer.Answer(&answer)
			if err != nil {
				ll.E(err).St().Err()
				return
			}
		}()

	})

	peer.OnIceCandidate(func(candidate *webrtc.ICECandidateInit, target int) {

		ll := l.Clone().F(log.FF{"peer": peer.GetUsername(), "event": "on-ice-candidate"}).Dbg().TrcF("candidate=%v", *candidate)

		if candidate == nil {
			ll.Dbg("done")
			return
		}

		err := avpPeer.AddICECandidate(*candidate, target)
		if err != nil {
			ll.E(err).St().Err()
		}

	})

	sdp, err := peer.Join(ctx, roomId, &offer)
	if err != nil {
		return nil, err
	}

	avpPeer.OnICECandidate(func(candidate *webrtc.ICECandidate, target int) {

		ll := l.Clone().F(log.FF{"peer": "avp", "event": "on-ice-candidate"}).Dbg()

		if candidate == nil {
			ll.Dbg("done")
			// gathering done
			return
		}

		ll.TrcF("candidate=%v", *candidate)

		err := peer.Trickle(candidate.ToJSON(), target)
		if err != nil {
			ll.E(err).St().Err()
		}
	})

	err = avpPeer.SetRemoteDescription(*sdp)
	if err != nil {
		return nil, err
	}

	return &webmSaverRoomRecorderImpl{
		cancelTransportFn: cancel,
	}, nil
}

func (r *webmSaverRoomRecorderImpl) Close() {
	r.cancelTransportFn()
}
