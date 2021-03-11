package storage

import (
	"context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/logger"
	"time"
)

type room struct {
	kitStorage.BaseDto
	Id       string     `gorm:"column:id"`
	Details  string     `gorm:"column:details"`
	OpenedAt *time.Time `gorm:"column:opened_at"`
	ClosedAt *time.Time `gorm:"column:closed_at"`
}

type storageImpl struct {
	c *container
}

func newStorage(c *container) *storageImpl {
	s := &storageImpl{c}
	return s
}

func (s *storageImpl) l() log.CLogger {
	return logger.L().Cmp("room-storage")
}

func (s *storageImpl) Create(ctx context.Context, room *domain.Room) (*domain.Room, error) {

	l := s.l().C(ctx).Mth("create").F(log.FF{"id": room.Id})

	dto := s.toRoomDto(room)

	t := time.Now().UTC()
	dto.CreatedAt, dto.UpdatedAt = t, t

	result := s.c.Db.Instance.Create(dto)
	if result.Error != nil {
		l.E(result.Error).St().Err()
		return nil, result.Error
	} else {
		l.Dbg("created")
		return s.toRoomDomain(dto), nil
	}
}

func (s *storageImpl) Update(ctx context.Context, room *domain.Room) (*domain.Room, error) {

	l := s.l().C(ctx).Mth("update").F(log.FF{"id": room.Id})

	dto := s.toRoomDto(room)

	dto.UpdatedAt = time.Now().UTC()

	result := s.c.Db.Instance.Save(dto)
	if result.Error != nil {
		l.E(result.Error).St().Err()
		return nil, result.Error
	} else {
		l.Dbg("updated")
		return s.toRoomDomain(dto), nil
	}
}

func (s *storageImpl) Get(ctx context.Context, roomId string) *domain.Room {
	dto := &room{Id: roomId}
	s.c.Db.Instance.First(dto)
	return s.toRoomDomain(dto)
}
