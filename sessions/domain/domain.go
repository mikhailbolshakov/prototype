package domain

import (
	"context"
	"time"
)

type LoginRequest struct {
	Username  string
	Password  string
	ChatLogin bool
}

type LoginResponse struct {
	SessionId string
}

type LogoutRequest struct {
	UserId string
}

type Session struct {
	Id            string
	UserId        string
	Username      string
	ChatUserId    string
	ChatSessionId string
	LoginAt       time.Time
}

type GetByUserRequest struct {
	UserId   string
	Username string
}

type SessionsService interface {
	Login(ctx context.Context, rq *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, rq *LogoutRequest) error
	Get(ctx context.Context, sid string) (*Session, error)
	GetByUser(ctx context.Context, rq *GetByUserRequest) ([]*Session, error)
	AuthSession(ctx context.Context, sid string) (*Session, error)
}

type SessionInfo struct {
	Id             string
	StartAt        time.Time
	SentWsMessages uint32
	ChatSessionId  string
}

type UserSessionInfo struct {
	UserId     string
	ChatUserId string
	Sessions   []*SessionInfo
}

type TotalSessionInfo struct {
	ActiveCount      int
	ActiveUsersCount int
}

type SessionMonitor interface {
	GetUserSessions(ctx context.Context, userId string) *UserSessionInfo
	GetTotalSessions(ctx context.Context) *TotalSessionInfo
}
