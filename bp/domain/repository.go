package domain

import (
	"github.com/Nerzal/gocloak/v7"
	kit "gitlab.medzdrav.ru/prototype/kit/config"
	pbChat "gitlab.medzdrav.ru/prototype/proto/chat"
	pbServ "gitlab.medzdrav.ru/prototype/proto/services"
	pbTask "gitlab.medzdrav.ru/prototype/proto/tasks"
	pbUser "gitlab.medzdrav.ru/prototype/proto/users"
	"time"
)

type ChatService interface {
	CreateClientChannel(rq *pbChat.CreateClientChannelRequest) (string, error)
	GetChannelsForUserAndExpert(userId, expertId string) ([]string, error)
	Subscribe(userId, channelId string) error
	CreateUser(rq *pbChat.CreateUserRequest) (string, error)
	DeleteUser(userId string) error
	AskBot(rq *pbChat.AskBotRequest) (*pbChat.AskBotResponse, error)
	Post(message, channelId, userId string, ephemeral, fromBot bool) error
	PredefinedPost(channelId, userId, code string, ephemeral, fromBot bool, params map[string]interface{}) error
}

type ConfigService interface {
	Get() (*kit.Config, error)
}

type BalanceService interface {
	Add(rq *pbServ.ChangeServicesRequest) (*pbServ.UserBalance, error)
	GetBalance(rq *pbServ.GetBalanceRequest) (*pbServ.UserBalance, error)
	WriteOff(rq *pbServ.ChangeServicesRequest) (*pbServ.UserBalance, error)
	Lock(rq *pbServ.ChangeServicesRequest) (*pbServ.UserBalance, error)
	CancelLock(rq *pbServ.ChangeServicesRequest) (*pbServ.UserBalance, error)
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
	Create(userId, serviceTypeId string, details map[string]interface{}) (*Delivery, error)
	GetDelivery(deliveryId string) (*Delivery, error)
	Cancel(deliveryId string, cancelTime *time.Time) (*Delivery, error)
	Complete(deliveryId string, completeTime *time.Time) (*Delivery, error)
	UpdateDetails(id string, details map[string]interface{}) (*Delivery, error)
}

type Status struct {
	Status    string
	SubStatus string
}

type Type struct {
	Type    string
	SubType string
}

type Assignee struct {
	Type     string
	Group    string
	UserId   string
	Username string
	At       *time.Time
}

type Reported struct {
	Type     string
	UserId   string
	Username string
	At       *time.Time
}

type Task struct {
	Id          string
	Num         string
	Type        *Type
	Status      *Status
	Reported    *Reported
	DueDate     *time.Time
	Assignee    *Assignee
	Description string
	Title       string
	Details     string
	ChannelId   string
}

type TaskService interface {
	GetByChannelId(channelId string) []*pbTask.Task
	New(rq *pbTask.NewTaskRequest) (*pbTask.Task, error)
	MakeTransition(rq *pbTask.MakeTransitionRequest) error
	Search(rq *pbTask.SearchRequest) ([]*pbTask.Task, error)
}

type UserService interface {
	Get(id string) *pbUser.User
	Activate(userId string) (*pbUser.User, error)
	Delete(userId string) (*pbUser.User, error)
	SetClientDetails(userId string, details *pbUser.ClientDetails) (*pbUser.User, error)
	SetMMUserId(userId, mmId string) (*pbUser.User, error)
	SetKKUserId(userId, kkId string) (*pbUser.User, error)
	GetByMMId(mmUserId string) (*pbUser.User, error)
}

type KeycloakProvider func() gocloak.GoCloak
