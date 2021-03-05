package domain

import (
	"context"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	sessionPb "gitlab.medzdrav.ru/prototype/proto/sessions"
)

type WebrtcStorage interface {
}


type ConfigService interface {
	Get(ctx context.Context) (*kit.Config, error)
}

type SessionsService interface {
	AuthSession(ctx context.Context, sid string) (*sessionPb.Session, error)
}

type RoomCoordinator interface {
	GetOrCreate(ctx context.Context, meta *RoomMeta) (bool, error)
	Close(ctx context.Context, roomId string)
}

type Sfu interface {
	
}