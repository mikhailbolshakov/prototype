package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"time"
)

const (
	TT_CLIENT       = "client"
	TST_MED_REQUEST = "medical-request"
	TST_EXPERT_CONSULTATION = "expert-consultation"

	TS_EMPTY  = "#"
	TS_OPEN   = "open"
	TS_CLOSED = "closed"

	TSS_EMPTY                    = "#"
	TSS_REPORTED      = "reported"
	TSS_ON_ASSIGNMENT = "on-assignment"
	TSS_ASSIGNED      = "assigned"
	TSS_IN_PROGRESS   = "in-progress"
	TSS_ON_HOLD       = "on-hold"
	TSS_CANCELLED     = "cancelled"
	TSS_SOLVED        = "solved"

	G_CLIENT     = "client"
	G_CONSULTANT = "consultant"
	G_EXPERT     = "expert"

	NUM_GEN_TYPE_RANDOM = "random"
	NUM_GEN_TYPE_SEQ    = "sequence"
)

type Status struct {
	Status    string
	SubStatus string
}

type Type struct {
	Type    string
	SubType string
}

type Assignee struct {
	Group string
	User  string
	At    *time.Time
}

type Reported struct {
	By string
	At *time.Time
}

type Transition struct {
	// transition id (must be unique over status model)
	Id string
	// source status
	From *Status
	// target status
	To *Status
	// for those groups the transition is allowed
	AllowAssignGroups []string
	// if not empty the task is assigned onto the group once transition happens
	AutoAssignGroup string
	// if true its an initial transition that is applied on creation a new task
	// there must be one and the only one transition with Initial flag
	Initial bool
	// if a topic is specified task is sent to Queue
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
	History     []*History
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
	Status   *Status
	Assignee *Assignee
	Type     *Type
	Num      string
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
