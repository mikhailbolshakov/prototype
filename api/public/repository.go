package public

import (
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	bpPb "gitlab.medzdrav.ru/prototype/proto/bp"
	servPb "gitlab.medzdrav.ru/prototype/proto/services"
	taskPb "gitlab.medzdrav.ru/prototype/proto/tasks"
	userPb "gitlab.medzdrav.ru/prototype/proto/users"
	"time"
)

type BpService interface {
	StartProcess(rq *bpPb.StartProcessRequest) (*bpPb.StartProcessResponse, error)
}

type ConfigService interface {
	Get() (*kit.Config, error)
}

type BalanceService interface {
	Add(rq *servPb.ChangeServicesRequest) (*servPb.UserBalance, error)
	GetBalance(rq *servPb.GetBalanceRequest) (*servPb.UserBalance, error)
	WriteOff(rq *servPb.ChangeServicesRequest) (*servPb.UserBalance, error)
	Lock(rq *servPb.ChangeServicesRequest) (*servPb.UserBalance, error)
	CancelLock(rq *servPb.ChangeServicesRequest) (*servPb.UserBalance, error)
}

type DeliveryService interface {
	Create(userId, serviceTypeId string, details map[string]interface{}) (*servPb.Delivery, error)
	GetDelivery(deliveryId string) (*servPb.Delivery, error)
	Cancel(deliveryId string, cancelTime *time.Time) (*servPb.Delivery, error)
	Complete(deliveryId string, completeTime *time.Time) (*servPb.Delivery, error)
	UpdateDetails(id string, details map[string]interface{}) (*servPb.Delivery, error)
}

type TaskService interface {
	New(rq *taskPb.NewTaskRequest) (*taskPb.Task, error)
	MakeTransition(rq *taskPb.MakeTransitionRequest) (*taskPb.Task, error)
	SetAssignee(rq *taskPb.SetAssigneeRequest) (*taskPb.Task, error)
	GetById(id string) (*taskPb.Task, error)
	Search(rq *taskPb.SearchRequest) (*taskPb.SearchResponse, error)
	GetAssignmentLog(rq *taskPb.AssignmentLogRequest) (*taskPb.AssignmentLogResponse, error)
	GetHistory(taskId string) (*taskPb.GetHistoryResponse, error)
}

type UserService interface {
	Get(id string) *userPb.User
	CreateClient(request *userPb.CreateClientRequest) (*userPb.User, error)
	CreateConsultant(request *userPb.CreateConsultantRequest) (*userPb.User, error)
	CreateExpert(request *userPb.CreateExpertRequest) (*userPb.User, error)
	Search(request *userPb.SearchRequest) (*userPb.SearchResponse, error)
}
