package storage

import "time"

type BaseDto struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
