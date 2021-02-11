package domain

import (
	"context"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"time"
)

const (
	TT_CLIENT                = "client"
	TST_REQUEST              = "common-request"
	TST_DENTIST_CONSULTATION = "dentist-consultation"
	TST_CLIENT_FEEDBACK      = "client-feedback"
	TST_MED_REQUEST          = "med-request"
	TST_LAWYER_REQUEST       = "lawyer-request"

	TT_TST  = "test"
	TST_TST = "test"

	TS_EMPTY  = "#"
	TS_OPEN   = "open"
	TS_CLOSED = "closed"

	TSS_EMPTY         = "#"
	TSS_REPORTED      = "reported"
	TSS_ON_ASSIGNMENT = "on-assignment"
	TSS_ASSIGNED      = "assigned"
	TSS_IN_PROGRESS   = "in-progress"
	TSS_ON_HOLD       = "on-hold"
	TSS_CANCELLED     = "cancelled"
	TSS_SOLVED        = "solved"

	USR_TYPE_CLIENT     = "client"
	USR_TYPE_CONSULTANT = "consultant"
	USR_TYPE_EXPERT     = "expert"

	USR_GRP_CLIENT            = "client"
	USR_GRP_CONSULTANT_LAWYER = "consultant-lawyer"
	USR_GRP_CONSULTANT_MED    = "consultant-med"
	USR_GRP_CONSULTANT_COMMON = "consultant"
	USR_GRP_DOCTOR_DENTIST    = "doctor-dentist"

	NUM_GEN_TYPE_RANDOM = "random"
	NUM_GEN_TYPE_SEQ    = "sequence"
)

type Status struct {
	Status    string
	SubStatus string
}

func (s *Status) Equals(another *Status) bool {
	if another == nil {
		return false
	}

	return s.Status == another.Status && s.SubStatus == another.SubStatus
}

type Type struct {
	Type    string
	SubType string
}

func (t *Type) Equals(another *Type) bool {
	if another == nil {
		return false
	}

	return t.Type == another.Type && t.SubType == another.SubType
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

type Transition struct {
	// transition id (must be unique over status model)
	Id string
	// if true its an initial transition that is applied on creation a new task
	// there must be one and the only one transition with Initial flag
	Initial bool
	// source status
	From *Status
	// target status
	To *Status
	// if not empty the task is assigned onto the group once transition happens
	AutoAssignType string
	// if not empty the task is assigned onto the group once transition happens
	AutoAssignGroup string
	// on making a transition if an assigned user isn't set, error occurs
	AssignedUserMandatory bool
	// if specified, a message will be send to this queue after a transition successfully has been made
	QueueTopic string
}

type StatusModel struct {
	Transitions []*Transition
}

type Attributes struct {
	AllowSchedule     bool
	AllowNotification bool
}

type NumGenerationRule struct {
	Prefix         string
	GenerationType string
}

type AssignmentSource struct {
	Status   *Status
	Assignee *Assignee
}

type AssignmentTarget struct {
	Status *Status
}

type UserPool struct {
	Type     string
	Group    string
	Statuses []string
}

type AssignmentRule struct {
	Code                  string
	Description           string
	DistributionAlgorithm string
	UserPool              *UserPool
	Source                *AssignmentSource
	Target                *AssignmentTarget
}

type Config struct {
	Id              string
	Type            *Type
	NumGenRule      *NumGenerationRule
	StatusModel     *StatusModel
	AssignmentRules []*AssignmentRule
}

type TimeUnit string

type BeforeDueDate struct {
	Unit  TimeUnit
	Value uint
}

type SpecificTime struct {
	At *time.Time
}

type Reminder struct {
	BeforeDueDate *BeforeDueDate
	SpecificTime  *SpecificTime
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
	Details     map[string]interface{}
	ChannelId   string
	Reminders   []*Reminder
}

type History struct {
	Id        string
	TaskId    string
	Status    *Status
	Assignee  *Assignee
	ChangedBy string
	ChangedAt time.Time
}

type SearchCriteria struct {
	*common.PagingRequest
	Status    *Status
	Assignee  *Assignee
	Type      *Type
	Num       string
	ChannelId string
}

type SearchResponse struct {
	*common.PagingResponse
	Tasks []*Task
}

type AssignmentLog struct {
	Id              string
	StartTime       time.Time
	FinishTime      *time.Time
	Status          string
	RuleCode        string
	RuleDescription string
	UsersInPool     int
	TasksToAssign   int
	Assigned        int
	Error           string
}

type AssignmentLogCriteria struct {
	*common.PagingRequest
	StartTimeAfter  *time.Time
	StartTimeBefore *time.Time
}

type AssignmentLogResponse struct {
	*common.PagingResponse
	Logs []*AssignmentLog
}

type ConfigService interface {
	GetAll(ctx context.Context) []*Config
	// get whole configuration
	Get(ctx context.Context, t *Type) (*Config, error)
	// if status is final
	IsFinalStatus(ctx context.Context, t *Type, s *Status) bool
	// retrieves a list of available statuses which might be set for the task
	NextTransitions(ctx context.Context, t *Type, currentStatus *Status) ([]*Transition, error)
	// get initial transition
	InitialTransition(ctx context.Context, t *Type) (*Transition, error)
	// get transition by source/target statuses
	FindTransition(ctx context.Context, t *Type, current, target *Status) (*Transition, error)
}

type TaskService interface {
	// register a new task
	New(ctx context.Context, task *Task) (*Task, error)
	// set status
	MakeTransition(ctx context.Context, taskId, transitionId string) (*Task, error)
	// assign
	SetAssignee(ctx context.Context, taskId string, target *Assignee) (*Task, error)
	// get by Id
	Get(ctx context.Context, taskId string) *Task
	// get tasks by channel
	GetByChannel(ctx context.Context, channelId string) []*Task
	// update task
	Update(ctx context.Context, task *Task) (*Task, error)
	// get assignment tasks execution log
	GetAssignmentLog(ctx context.Context, cr *AssignmentLogCriteria) (*AssignmentLogResponse, error)
	// get task history
	GetHistory(ctx context.Context, taskId string) []*History
	// search tasks by the criteria
	Search(ctx context.Context, cr *SearchCriteria) (*SearchResponse, error)
}

// Daemon assigns tasks to users according to rules
type AssignmentDaemon interface {
	// run daemon
	Run(ctx context.Context)
	// stop daemon
	Stop(ctx context.Context) error
	// init daemon
	Init(ctx context.Context) error
}

type TaskSchedulerHandler func(ctx context.Context, taskId string)

type TaskScheduler interface {
	SetReminderHandler(h TaskSchedulerHandler)
	SetDueDateHandler(h TaskSchedulerHandler)
	StartAsync(ctx context.Context, )
	ScheduleTask(ctx context.Context, ts *Task)
}
