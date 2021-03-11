package domain

import (
	"context"
	"gitlab.medzdrav.ru/prototype/proto/config"
	sessionPb "gitlab.medzdrav.ru/prototype/proto/sessions"
)

type RoomStorage interface {
	Create(ctx context.Context, room *Room) (*Room, error)
	Update(ctx context.Context, room *Room) (*Room, error)
	Get(ctx context.Context, roomId string) *Room
}

type ConfigService interface {
	Get(ctx context.Context) (*config.Config, error)
}

type SessionsService interface {
	AuthSession(ctx context.Context, sid string) (*sessionPb.Session, error)
}

type RoomCoordinator interface {
	GetOrCreate(ctx context.Context, meta *RoomMeta) (bool, error)
	Close(ctx context.Context, roomId string)
}
