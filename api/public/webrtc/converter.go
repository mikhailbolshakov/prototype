package webrtc

import (
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/webrtc"
)

func (c *ctrlImpl) toRoomApi(r *pb.Room) *Room {

	if r == nil {
		return nil
	}

	res := &Room{
		Id:       r.Id,
		Details:  &RoomDetails{
			ChannelId:    r.Details.ChannelId,
			Participants: []*RoomParticipants{},
		},
		OpenedAt: grpc.PbTSToTime(r.OpenedAt),
		ClosedAt: grpc.PbTSToTime(r.ClosedAt),
	}

	for _, p := range r.Details.Participants {
		res.Details.Participants = append(res.Details.Participants, &RoomParticipants{
			UserId:   p.UserId,
			Username: p.Username,
			JoinedAt: *(grpc.PbTSToTime(p.JoinedAt)),
			LeaveAt:  grpc.PbTSToTime(p.LeaveAt),
		})
	}

	return res
}
