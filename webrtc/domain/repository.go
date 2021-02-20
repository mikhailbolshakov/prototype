package domain

import "context"

// RoomStorage is responsible for rooms persistence
type RoomStorage interface {
	Save(ctx context.Context, r *Room) error
	SaveAsync(ctx context.Context, r *Room)
}
