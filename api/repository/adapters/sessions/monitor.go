package sessions

import (
	"context"
	pb "gitlab.medzdrav.ru/prototype/proto/sessions"
)

type monitorImpl struct {
	pb.MonitorClient
}

func newMonitorImpl() *monitorImpl {
	a := &monitorImpl{}
	return a
}

func (m *monitorImpl) UserSessions(ctx context.Context, userId string) (*pb.UserSessionsInfo, error) {
	return m.MonitorClient.UserSessions(ctx, &pb.UserSessionsRequest{UserId: userId})
}

func (m *monitorImpl) TotalSessions(ctx context.Context) (*pb.TotalSessionInfo, error) {
	return m.MonitorClient.TotalSessions(ctx, &pb.SessionsTotalRequest{})
}