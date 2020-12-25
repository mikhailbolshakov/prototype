package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"time"
)

type Status struct {
	Status string
	SubStatus string
}

type Type struct {
	Type string
	SubType string
}

type Assignee struct {
	Group string
	User string
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
}

type StatusModel struct {
	Transitions []*Transition
}

type Attributes struct {
	AllowSchedule bool
	AllowNotification bool
}

const (
	NUM_GEN_TYPE_RANDOM = "random"
	NUM_GEN_TYPE_SEQ = "sequence"
)

type NumGenerationRule struct {
	Prefix string
	GenerationType string
}

type Config struct {
	Id string
	Type *Type
	NumGenRule *NumGenerationRule
	StatusModel *StatusModel
}

type Task struct {
	Id string
	Num string
	Type *Type
	Status *Status
	ReportedBy string
	ReportedAt *time.Time
	DueDate *time.Time
	Assignee *Assignee
	Description string
	Title string
	Details string
	History []*HistoryItem
}

type HistoryItem struct {
	Id string
	TaskId string
	Status string
	SubStatus string
	AssignedGroup string
	AssignedUser string
	ChangedBy string
	ChangedAt time.Time
}

type SearchCriteria struct {
	*common.PagingRequest
}

type SearchResponse struct {
	*common.PagingResponse
	Tasks []*Task
}
