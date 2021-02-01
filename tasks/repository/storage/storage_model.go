package storage

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"time"
)

type Task struct {
	kit.BaseDto
	Id               string     `gorm:"column:id"`
	Num              string     `gorm:"column:num"`
	Type             string     `gorm:"column:type"`
	SubType          string     `gorm:"column:subtype"`
	Status           string     `gorm:"column:status"`
	SubStatus        string     `gorm:"column:substatus"`
	ReportedType     string     `gorm:"column:reported_type"`
	ReportedUserId   string     `gorm:"column:reported_user_id"`
	ReportedUsername string     `gorm:"column:reported_username"`
	ReportedAt       time.Time  `gorm:"column:reported_at"`
	DueDate          *time.Time `gorm:"column:due_date"`
	AssigneeType     string     `gorm:"column:assignee_type"`
	AssigneeGroup    string     `gorm:"column:assignee_group"`
	AssigneeUserId   string     `gorm:"column:assignee_user_id"`
	AssigneeUsername string     `gorm:"column:assignee_username"`
	AssigneeAt       *time.Time `gorm:"column:assignee_at"`
	Description      string     `gorm:"column:description"`
	Title            string     `gorm:"column:title"`
	Details          string     `gorm:"column:details"`
	Reminders        string     `gorm:"column:reminders"`
	ChannelId        string     `gorm:"column:channel_id"`
}

type History struct {
	Id               string     `gorm:"column:id"`
	TaskId           string     `gorm:"column:task_id"`
	Status           string     `gorm:"column:status"`
	SubStatus        string     `gorm:"column:substatus"`
	AssigneeType     string     `gorm:"column:assignee_type"`
	AssigneeGroup    string     `gorm:"column:assignee_group"`
	AssigneeUserId   string     `gorm:"column:assignee_user_id"`
	AssigneeUsername string     `gorm:"column:assignee_username"`
	AssigneeAt       *time.Time `gorm:"column:assignee_at"`
	ChangedBy        string     `gorm:"column:changed_by"`
	ChangedAt        time.Time  `gorm:"column:changed_at"`
}

type SearchCriteria struct {
	*common.PagingRequest
	Num              string `gorm:"column:num"`
	Status           string `gorm:"column:status"`
	SubStatus        string `gorm:"column:substatus"`
	AssigneeType     string `gorm:"column:assignee_type"`
	AssigneeGroup    string `gorm:"column:assignee_group"`
	AssigneeUserId   string `gorm:"column:assignee_user_id"`
	AssigneeUsername string `gorm:"column:assignee_username"`
	Type             string `gorm:"column:type"`
	SubType          string `gorm:"column:subtype"`
	ChannelId        string `gorm:"column:channel_id"`
}

type SearchResponse struct {
	*common.PagingResponse
	Tasks []*Task
}

type AssignmentLog struct {
	Id              string     `gorm:"column:id"`
	StartTime       time.Time  `gorm:"column:start_time"`
	FinishTime      *time.Time `gorm:"column:finish_time"`
	Status          string     `gorm:"column:status"`
	RuleCode        string     `gorm:"column:rule_code"`
	RuleDescription string     `gorm:"column:rule_description"`
	UsersInPool     int        `gorm:"column:users_in_pool"`
	TasksToAssign   int        `gorm:"column:tasks_to_assign"`
	Assigned        int        `gorm:"column:assigned"`
	Error           string     `gorm:"column:error"`
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
