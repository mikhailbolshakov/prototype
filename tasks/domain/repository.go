package domain

import (
	"context"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type TaskStorage interface {
	Create(ctx context.Context, t *Task) (*Task, error)
	Get(ctx context.Context, id string) *Task
	GetByIds(ctx context.Context, id []string) []*Task
	Update(ctx context.Context, t *Task) (*Task, error)
	GetByChannel(ctx context.Context, channelId string) []*Task
	CreateHistory(ctx context.Context, h *History) (*History, error)
	Search(ctx context.Context, cr *SearchCriteria) (*SearchResponse, error)
	SaveAssignmentLog(ctx context.Context, l *AssignmentLog) (*AssignmentLog, error)
	GetAssignmentLog(ctx context.Context, c *AssignmentLogCriteria) (*AssignmentLogResponse, error)
	GetHistory(ctx context.Context, taskId string) []*History
}

type ChatService interface {
	Post(ctx context.Context, message, channelId, userId string, ephemeral, fromBot bool) error
	PredefinedPost(ctx context.Context, channelId, userId, code string, ephemeral, fromBot bool, params map[string]interface{}) error
}

type CfgService interface {
	Get(ctx context.Context) (*kit.Config, error)
}

type UserService interface {
	Get(ctx context.Context, id, username string) *pb.User
	Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error)
}
