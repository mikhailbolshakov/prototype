package tasks

import "time"

const (
	QUEUE_TOPIC_TASK_SOLVED_STATUS = "tasks.solved"
	QUEUE_TOPIC_TASK_ASSIGN_STATUS = "tasks.assigned"
	QUEUE_TOPIC_TASK_DUEDATE = "tasks.duedate"
)

type TaskStatusPayload struct {
	Status    string
	SubStatus string
}

type TaskTypePayload struct {
	Type    string
	SubType string
}

type TaskAssigneePayload struct {
	Type     string
	Group    string
	UserId   string
	Username string
	At       *time.Time
}

type TaskReportedPayload struct {
	Type     string
	UserId   string
	Username string
	At       *time.Time
}

type TaskMessagePayload struct {
	Id          string
	Num         string
	Type        *TaskTypePayload
	Status      *TaskStatusPayload
	Reported    *TaskReportedPayload
	DueDate     *time.Time
	Assignee    *TaskAssigneePayload
	Description string
	Title       string
	Details     string
	ChannelId   string
}