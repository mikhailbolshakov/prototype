package grpc

import (
	"encoding/json"
	"github.com/pion/webrtc/v3"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	webrtcPb "gitlab.medzdrav.ru/prototype/proto/webrtc"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
)

func (s *Server) toSdpDomain(sdp []byte) (*webrtc.SessionDescription, error) {

	var offer *webrtc.SessionDescription
	err := json.Unmarshal(sdp, &offer)
	if err != nil {
		return nil, err
	}
	return offer, nil

}

func (s *Server) toRoomPb(r *domain.Room) *webrtcPb.Room {

	if r == nil {
		return nil
	}

	res := &webrtcPb.Room{
		Id:       r.Id,
		Details:  &webrtcPb.RoomDetails{
			ChannelId:    r.Details.ChannelId,
			Participants: []*webrtcPb.RoomParticipants{},
		},
		OpenedAt: grpc.TimeToPbTS(r.OpenedAt),
		ClosedAt: grpc.TimeToPbTS(r.ClosedAt),
	}

	for _, p := range r.Details.Participants {
		res.Details.Participants = append(res.Details.Participants, &webrtcPb.RoomParticipants{
			UserId:   p.UserId,
			Username: p.Username,
			JoinedAt: grpc.TimeToPbTS(&p.JoinedAt),
			LeaveAt:  grpc.TimeToPbTS(p.LeaveAt),
		})
	}

	return res
}

