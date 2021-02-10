package domain

import (
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
)

type TaskStorage interface {
	Create(t *Task) (*Task, error)
	Get(id string) *Task
	GetByIds(id []string) []*Task
	Update(t *Task) (*Task, error)
	GetByChannel(channelId string) []*Task
	CreateHistory(h *History) (*History, error)
	Search(cr *SearchCriteria) (*SearchResponse, error)
	SaveAssignmentLog(l *AssignmentLog) (*AssignmentLog, error)
	GetAssignmentLog(c *AssignmentLogCriteria) (*AssignmentLogResponse, error)
	GetHistory(taskId string) []*History
}

type ChatService interface {
	Post(message, channelId, userId string, ephemeral, fromBot bool) error
	PredefinedPost(channelId, userId, code string, ephemeral, fromBot bool, params map[string]interface{}) error
}

type CfgService interface {
	Get() (*kit.Config, error)
}

type UserService interface {
	Get(id, username string) *pb.User
	Search(request *pb.SearchRequest) (*pb.SearchResponse, error)
}
