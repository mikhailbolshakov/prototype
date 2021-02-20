package domain

import (
	"context"
	"github.com/Nerzal/gocloak/v7"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	pbChat "gitlab.medzdrav.ru/prototype/proto/chat"
	pbServ "gitlab.medzdrav.ru/prototype/proto/services"
	pbTask "gitlab.medzdrav.ru/prototype/proto/tasks"
	pbUser "gitlab.medzdrav.ru/prototype/proto/users"
	"time"
)

type ChatService interface {
	CreateClientChannel(ctx context.Context, rq *pbChat.CreateClientChannelRequest) (string, error)
	GetChannelsForUserAndExpert(ctx context.Context, userId, expertId string) ([]string, error)
	Subscribe(ctx context.Context, userId, channelId string) error
	CreateUser(ctx context.Context, rq *pbChat.CreateUserRequest) (string, error)
	DeleteUser(ctx context.Context, userId string) error
	AskBot(ctx context.Context, rq *pbChat.AskBotRequest) (*pbChat.AskBotResponse, error)
	Post(ctx context.Context, message, channelId, userId string, ephemeral bool) error
	PredefinedPost(ctx context.Context, channelId, userId, code string, ephemeral bool, params map[string]interface{}) error
}

type ConfigService interface {
	Get() (*kit.Config, error)
}

type BalanceService interface {
	Add(ctx context.Context, rq *pbServ.ChangeServicesRequest) (*pbServ.UserBalance, error)
	GetBalance(ctx context.Context, rq *pbServ.GetBalanceRequest) (*pbServ.UserBalance, error)
	WriteOff(ctx context.Context, rq *pbServ.ChangeServicesRequest) (*pbServ.UserBalance, error)
	Lock(ctx context.Context, rq *pbServ.ChangeServicesRequest) (*pbServ.UserBalance, error)
	CancelLock(ctx context.Context, rq *pbServ.ChangeServicesRequest) (*pbServ.UserBalance, error)
}

type Delivery struct {
	Id string
	UserId string
	ServiceTypeId string
	Status string
	StartTime time.Time
	FinishTime *time.Time
	Details map[string]interface{}
}

type DeliveryService interface {
	Create(ctx context.Context, userId, serviceTypeId string, details map[string]interface{}) (*Delivery, error)
	GetDelivery(ctx context.Context, deliveryId string) (*Delivery, error)
	Cancel(ctx context.Context, deliveryId string, cancelTime *time.Time) (*Delivery, error)
	Complete(ctx context.Context, deliveryId string, completeTime *time.Time) (*Delivery, error)
	UpdateDetails(ctx context.Context, id string, details map[string]interface{}) (*Delivery, error)
}

type TaskService interface {
	GetByChannelId(ctx context.Context, channelId string) []*pbTask.Task
	New(ctx context.Context, rq *pbTask.NewTaskRequest) (*pbTask.Task, error)
	MakeTransition(ctx context.Context, rq *pbTask.MakeTransitionRequest) error
	Search(ctx context.Context, rq *pbTask.SearchRequest) ([]*pbTask.Task, error)
}

type UserService interface {
	Get(ctx context.Context, id string) *pbUser.User
	Activate(ctx context.Context, userId string) (*pbUser.User, error)
	Delete(ctx context.Context, userId string) (*pbUser.User, error)
	SetClientDetails(ctx context.Context, userId string, details *pbUser.ClientDetails) (*pbUser.User, error)
	SetMMUserId(ctx context.Context, userId, mmId string) (*pbUser.User, error)
	SetKKUserId(ctx context.Context, userId, kkId string) (*pbUser.User, error)
	GetByMMId(ctx context.Context, mmUserId string) (*pbUser.User, error)
}

type KeycloakProvider func() gocloak.GoCloak
