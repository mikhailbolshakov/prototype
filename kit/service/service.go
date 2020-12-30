package service

type Service interface {
	Init() error
	Listen() error
	ListenAsync() error
}
