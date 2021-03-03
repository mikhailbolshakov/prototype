package impl

import (
	"context"
	"gitlab.medzdrav.ru/prototype/sessions/domain"
)

type monitorImpl struct {
	hub     Hub
	metrics domain.Metrics
}

func NewMonitorService(h Hub, metrics domain.Metrics) domain.SessionMonitor {
	return &monitorImpl{
		hub:     h,
		metrics: metrics,
	}
}

func (m *monitorImpl) GetUserSessions(ctx context.Context, userId string) *domain.UserSessionInfo {

	rs := &domain.UserSessionInfo{
		UserId:   userId,
		Sessions: []*domain.SessionInfo{},
	}

	for i, s := range m.hub.getByUserId(userId) {

		if i == 0 {
			rs.ChatUserId = s.getChatUserId()
		}

		rs.Sessions = append(rs.Sessions, &domain.SessionInfo{
			Id:             s.getId(),
			StartAt:        s.getStartAt(),
			SentWsMessages: s.getSentWsMessages(),
			ChatSessionId:  s.getChatSessionId(),
		})
	}

	return rs
}

func (m *monitorImpl) GetTotalSessions(ctx context.Context) *domain.TotalSessionInfo {

	return &domain.TotalSessionInfo{
		ActiveCount:      m.metrics.SessionsCount(),
		ActiveUsersCount: m.metrics.ConnectedUsersCount(),
	}

}
