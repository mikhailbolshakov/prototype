package storage

import (
	"encoding/json"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
)

func (s *storageImpl) toRoomDto(r *domain.Room) *room {

	if r == nil {
		return nil
	}

	det, _ :=  json.Marshal(r.Details)

	return &room{
		BaseDto:  kitStorage.BaseDto{
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.ModifiedAt,
			DeletedAt: r.DeletedAt,
		},
		Id:       r.Id,
		Details:  string(det),
		OpenedAt: r.OpenedAt,
		ClosedAt: r.ClosedAt,
	}

}

func (s *storageImpl) toRoomDomain(r *room) *domain.Room {

	if r == nil {
		return nil
	}

	details := &domain.RoomDetails{}
	if r.Details != "" {
		_ = json.Unmarshal([]byte(r.Details), &details)
	}

	return &domain.Room{
		Id:         r.Id,
		OpenedAt:   r.OpenedAt,
		ClosedAt:   r.ClosedAt,
		Details:    details,
		CreatedAt:  r.CreatedAt,
		ModifiedAt: r.UpdatedAt,
		DeletedAt:  r.DeletedAt,
	}

}
