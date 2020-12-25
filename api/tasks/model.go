package tasks

import "time"

type Status struct {
	Status string `json:"status"`
	SubStatus string `json:"substatus"`
}

type Type struct {
	Type string `json:"type"`
	SubType string `json:"subtype"`
}

type Assignee struct {
	Group string `json:"group"`
	User string `json:"user"`
	At *time.Time `json:"at"`
}

type Task struct {
	Id string `json:"id"`
	Num string `json:"num"`
	Type *Type `json:"type"`
	Status *Status `json:"status"`
	ReportedBy string  `json:"reportedBy"`
	ReportedAt *time.Time `json:"reportedAt"`
	DueDate *time.Time `json:"dueDate"`
	Assignee *Assignee `json:"assignee"`
	Description string `json:"description"`
	Title string `json:"title"`
	Details string `json:"details"`
}

type NewTaskRequest struct {
	Type *Type `json:"type"`
	ReportedBy string `json:"reportedBy"`
	ReportedAt *time.Time `json:"reportedAt"`
	DueDate *time.Time `json:"dueDate"`
	Assignee *Assignee `json:"assignee"`
	Description string `json:"description"`
	Title string `json:"title"`
}

