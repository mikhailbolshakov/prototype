package public

import (
	"context"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	bpPb "gitlab.medzdrav.ru/prototype/proto/bp"
	servPb "gitlab.medzdrav.ru/prototype/proto/services"
	taskPb "gitlab.medzdrav.ru/prototype/proto/tasks"
	userPb "gitlab.medzdrav.ru/prototype/proto/users"
	"time"
)

type BpService interface {
	StartProcess(ctx context.Context, rq *bpPb.StartProcessRequest) (*bpPb.StartProcessResponse, error)
}

type ConfigService interface {
	Get() (*kit.Config, error)
}

type BalanceService interface {
	Add(ctx context.Context, rq *servPb.ChangeServicesRequest) (*servPb.UserBalance, error)
	GetBalance(ctx context.Context, rq *servPb.GetBalanceRequest) (*servPb.UserBalance, error)
	WriteOff(ctx context.Context, rq *servPb.ChangeServicesRequest) (*servPb.UserBalance, error)
	Lock(ctx context.Context, rq *servPb.ChangeServicesRequest) (*servPb.UserBalance, error)
	CancelLock(ctx context.Context, rq *servPb.ChangeServicesRequest) (*servPb.UserBalance, error)
}

type DeliveryService interface {
	Create(ctx context.Context, userId, serviceTypeId string, details map[string]interface{}) (*servPb.Delivery, error)
	GetDelivery(ctx context.Context, deliveryId string) (*servPb.Delivery, error)
	Cancel(ctx context.Context, deliveryId string, cancelTime *time.Time) (*servPb.Delivery, error)
	Complete(ctx context.Context, deliveryId string, completeTime *time.Time) (*servPb.Delivery, error)
	UpdateDetails(ctx context.Context, id string, details map[string]interface{}) (*servPb.Delivery, error)
}

type TaskService interface {
	New(ctx context.Context, rq *taskPb.NewTaskRequest) (*taskPb.Task, error)
	MakeTransition(ctx context.Context, rq *taskPb.MakeTransitionRequest) (*taskPb.Task, error)
	SetAssignee(ctx context.Context, rq *taskPb.SetAssigneeRequest) (*taskPb.Task, error)
	GetById(ctx context.Context, id string) (*taskPb.Task, error)
	Search(ctx context.Context, rq *taskPb.SearchRequest) (*taskPb.SearchResponse, error)
	GetAssignmentLog(ctx context.Context, rq *taskPb.AssignmentLogRequest) (*taskPb.AssignmentLogResponse, error)
	GetHistory(ctx context.Context, taskId string) (*taskPb.GetHistoryResponse, error)
}

type UserService interface {
	Get(ctx context.Context, id string) *userPb.User
	CreateClient(ctx context.Context, request *userPb.CreateClientRequest) (*userPb.User, error)
	CreateConsultant(ctx context.Context, request *userPb.CreateConsultantRequest) (*userPb.User, error)
	CreateExpert(ctx context.Context, request *userPb.CreateExpertRequest) (*userPb.User, error)
	Search(ctx context.Context, request *userPb.SearchRequest) (*userPb.SearchResponse, error)
}
