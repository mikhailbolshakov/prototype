package monitoring

import "time"

type SessionInfo struct {
	Id             string    `json:"id"`
	StartAt        *time.Time `json:"startAt"`
	SentWsMessages uint32    `json:"sentWsMessages"`
	ChatSessionId  string    `json:"chatSessionId"`
}

type UserSessionInfo struct {
	UserId     string         `json:"userId"`
	ChatUserId string         `json:"chatUserId"`
	Sessions   []*SessionInfo `json:"sessions"`
}

type TotalSessionInfo struct {
	ActiveCount      int `json:"activeCount"`
	ActiveUsersCount int `json:"activeUsersCount"`
}
