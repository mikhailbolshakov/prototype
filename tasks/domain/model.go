package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"time"
)

const (
	TASK_TYPE_CLIENT         = "client"
	TASK_SUBTYPE_MED_REQUEST = "medical-request"

	TASK_STATUS_EMPTY  = "#"
	TASK_STATUS_OPEN   = "open"
	TASK_STATUS_CLOSED = "closed"

	TASK_SUBSTATUS_EMPTY         = "#"
	TASK_SUBSTATUS_REPORTED      = "reported"
	TASK_SUBSTATUS_ON_ASSIGNMENT = "on-assignment"
	TASK_SUBSTATUS_ASSIGNED      = "assigned"
	TASK_SUBSTATUS_IN_PROGRESS   = "in-progress"
	TASK_SUBSTATUS_ON_HOLD       = "on-hold"
	TASK_SUBSTATUS_CANCELLED     = "cancelled"
	TASK_SUBSTATUS_SOLVED        = "solved"

	GROUP_CLIENT     = "client"
	GROUP_CONSULTANT = "consultant"
	GROUP_EXPERT     = "expert"

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
