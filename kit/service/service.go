package service

type Service interface {
	Init() error
	ListenAsync() error
	Close()
}
