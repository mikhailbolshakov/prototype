package service

import "context"

type Service interface {
	Init(ctx context.Context) error
	ListenAsync(ctx context.Context) error
	Close(ctx context.Context)
}
