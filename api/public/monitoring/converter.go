package monitoring

import (
	"gitlab.medzdrav.ru/prototype/api/session"
)

func (c *ctrlImpl) toUserSessionsApi(si *session.UserSessionInfo) *UserSessionInfo {
	res := &UserSessionInfo{
		UserId:     si.UserId,
		ChatUserId: si.ChatUserId,
		Sessions: []*SessionInfo{},
	}

	for _, s := range si.Sessions {
		res.Sessions = append(res.Sessions, &SessionInfo{
			Id:             s.Id,
			StartAt:        s.StartAt,
			SentWsMessages: s.SentWsMessages,
			ChatSessionId:  s.ChatSessionId,
		})
	}

	return res
}

func (c *ctrlImpl) toTotalSessionsApi(si *session.TotalSessionInfo) *TotalSessionInfo {
	res := &TotalSessionInfo{
		ActiveCount:     si.ActiveCount,
		ActiveUsersCount: si.ActiveUsersCount,
	}
	return res
}
