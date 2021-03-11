package impl

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/logger"
	"time"
)

type roomServiceImpl struct {
	common.BaseService
	storage domain.RoomStorage
}

func NewRoomService(storage domain.RoomStorage) domain.RoomService {
	return &roomServiceImpl{
		storage: storage,
	}
}

func (r *roomServiceImpl) l() log.CLogger {
	return logger.L().Cmp("webrtc-room")
}

func (r *roomServiceImpl) Create(ctx context.Context, channelId string) (*domain.Room, error) {

	l := r.l().C(ctx).Mth("create")

	createdAt := time.Now().UTC()

	room := &domain.Room{
		Id: kit.NewId(),
		Details: &domain.RoomDetails{
			ChannelId:    channelId,
			Participants: []*domain.RoomParticipants{},
		},
		CreatedAt:  createdAt,
		ModifiedAt: createdAt,
	}

	room, err := r.storage.Create(ctx, room)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	l.F(log.FF{"room": room.Id}).Dbg()

	return room, nil

}

func (r *roomServiceImpl) Get(ctx context.Context, roomId string) *domain.Room {

	l := r.l().C(ctx).Mth("get").F(log.FF{"room": roomId})

	room := r.storage.Get(ctx, roomId)

	if room != nil {
		l.Trc("found")
	} else {
		l.Trc("not found")
	}

	return room
}

func (r *roomServiceImpl) Join(ctx context.Context, roomId, userId, username, peerId string) (*domain.Room, error) {

	l := r.l().C(ctx).Mth("join").F(log.FF{"room": roomId, "user": username, "peer": peerId})

	room := r.Get(ctx, roomId)
	if room == nil {
		err := fmt.Errorf("room not found")
		l.E(err).St().Err()
		return nil, err
	}

	if room.ClosedAt != nil {
		err := fmt.Errorf("room is closed")
		l.E(err).St().Err()
		return nil, err
	}

	// first participant is joining
	if len(room.Details.Participants) == 0 {
		t := time.Now().UTC()
		room.OpenedAt = &t
	}
	// search for participant with userId
	for _, p := range room.Details.Participants {
		if p.PeerId == peerId && p.LeaveAt == nil {
			err := fmt.Errorf("participant is in the room already")
			l.E(err).St().Err()
			return nil, err
		}
	}

	room.Details.Participants = append(room.Details.Participants, &domain.RoomParticipants{
		PeerId:   peerId,
		UserId:   userId,
		Username: username,
		JoinedAt: time.Now().UTC(),
	})

	room, err := r.storage.Update(ctx, room)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	l.Dbg("joined the room")

	return room, nil
}

func (r *roomServiceImpl) Leave(ctx context.Context, roomId, peerId string) (*domain.Room, error) {

	l := r.l().C(ctx).Mth("join").F(log.FF{"room": roomId, "peer": peerId})

	room := r.Get(ctx, roomId)
	if room == nil {
		err := fmt.Errorf("room not found")
		l.E(err).St().Err()
		return nil, err
	}

	// search for participant with userId
	leftTime := time.Now().UTC()
	found := false
	aliveExists := false
	for _, p := range room.Details.Participants {
		if p.PeerId == peerId && p.LeaveAt == nil {
			p.LeaveAt = &leftTime
			found = true
		}
		if p.LeaveAt == nil {
			aliveExists = true
		}
	}

	if !found {
		err := fmt.Errorf("participant not found")
		l.E(err).St().Warn()
	}

	if !aliveExists {
		room.ClosedAt = &leftTime
		l.Dbg("no alive participants, room is closed")
	}

	room, err := r.storage.Update(ctx, room)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	l.Dbg("left the room")

	return room, nil
}

func (r *roomServiceImpl) Close(ctx context.Context, roomId string) (*domain.Room, error) {

	l := r.l().C(ctx).Mth("close").F(log.FF{"room": roomId})

	room := r.Get(ctx, roomId)

	t := time.Now().UTC()
	room.ClosedAt = &t

	// enforce leaving all participants
	for _, p := range room.Details.Participants {
		if p.LeaveAt == nil {
			p.LeaveAt = &t
		}
	}

	room, err := r.storage.Update(ctx, room)
	if err != nil {
		l.E(err).St().Err()
		return nil, err
	}

	return room, nil
}
