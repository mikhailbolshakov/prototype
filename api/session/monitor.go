package session

import (
	"context"
	"time"
)

/*

Simple session monitor
it's just for dev purposes as it doesn't gather info for cluster configuration (when there are more than one hub instances)

!!! USER IT CAREFULLY
Since it locks the whole HUB for Read, so all new sessions are locked until monitor request is finished

*/

type SessionInfo struct {
	Id             string    `json:"id"`
	StartAt        time.Time `json:"startAt"`
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

type SessionMonitor interface {
	GetUserSessions(ctx context.Context, userId string) *UserSessionInfo
	GetTotalSessions(ctx context.Context) *TotalSessionInfo
}

func (h *hubImpl) GetUserSessions(ctx context.Context, userId string) *UserSessionInfo {
	h.RLock()
	defer h.RUnlock()

	rs := &UserSessionInfo{
		UserId:   userId,
		Sessions: []*SessionInfo{},
	}

	if sessions, ok := h.userSessions[userId]; ok {

		for i, s := range sessions {

			if i == 0 {
				rs.ChatUserId = s.getChatUserId()
			}

			rs.Sessions = append(rs.Sessions, &SessionInfo{
				Id:             s.getId(),
				StartAt:        s.getStartAt(),
				SentWsMessages: s.getSentWsMessages(),
				ChatSessionId:  s.getChatSessionId(),
			})
		}

	}
	return rs
}

func (h *hubImpl) GetTotalSessions(ctx context.Context) *TotalSessionInfo {
	h.RLock()
	defer h.RUnlock()

	rs := &TotalSessionInfo{}

	for _, s := range h.userSessions {
		rs.ActiveUsersCount++
		rs.ActiveCount += len(s)
	}

	return rs

}
