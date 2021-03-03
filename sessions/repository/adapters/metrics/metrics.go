package metrcics

import "go.uber.org/atomic"


//TODO: prometheus

type serviceImpl struct {
	sessions *atomic.Int32
	users    *atomic.Int32
}

func newImpl() *serviceImpl {
	m := &serviceImpl{
		sessions: atomic.NewInt32(0),
		users: atomic.NewInt32(0),
	}
	return m
}

func (m *serviceImpl) SessionsInc() {
	m.sessions.Inc()
}

func (m *serviceImpl) SessionsDec() {
	m.sessions.Dec()
}

func (m *serviceImpl) ConnectedUsersInc() {
	m.users.Inc()
}

func (m *serviceImpl) ConnectedUsersDec() {
	m.users.Dec()
}

func (m *serviceImpl) SessionsCount() int {
	return int(m.sessions.Load())
}

func (m *serviceImpl) ConnectedUsersCount() int {
	return int(m.users.Load())
}
