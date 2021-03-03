package monitoring

import (
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	sessionPb "gitlab.medzdrav.ru/prototype/proto/sessions"
)

func (c *ctrlImpl) toUserSessionsApi(si *sessionPb.UserSessionsInfo) *UserSessionInfo {
	res := &UserSessionInfo{
		UserId:     si.UserId,
		ChatUserId: si.ChatUserId,
		Sessions:   []*SessionInfo{},
	}

	for _, s := range si.Sessions {
		res.Sessions = append(res.Sessions, &SessionInfo{
			Id:             s.Id,
			StartAt:        grpc.PbTSToTime(s.StartAt),
			SentWsMessages: s.SentWsMessages,
			ChatSessionId:  s.ChatSessionId,
		})
	}

	return res
}

func (c *ctrlImpl) toTotalSessionsApi(si *sessionPb.TotalSessionInfo) *TotalSessionInfo {
	res := &TotalSessionInfo{
		ActiveCount:      int(si.ActiveCount),
		ActiveUsersCount: int(si.ActiveUsersCount),
	}
	return res
}
