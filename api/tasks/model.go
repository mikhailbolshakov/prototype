package tasks

import "time"

type Status struct {
	Status    string `json:"status"`
	SubStatus string `json:"substatus"`
}

type Type struct {
	Type    string `json:"type"`
	SubType string `json:"subtype"`
}

type Assignee struct {
	Group string     `json:"group,omitempty"`
	User  string     `json:"user,omitempty"`
	At    *time.Time `json:"at,omitempty"`
}

type Reported struct {
	By string     `json:"by"`
	At *time.Time `json:"at"`
}

type Task struct {
	Id          string     `json:"id"`
	Num         string     `json:"num"`
	Type        *Type      `json:"type"`
	Status      *Status    `json:"status"`
	Reported    *Reported  `json:"reported"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	Assignee    *Assignee  `json:"assignee"`
	Description string     `json:"description"`
	Title       string     `json:"title"`
	Details     string     `json:"details"`
}

type NewTaskRequest struct {
	Type        *Type      `json:"type"`
	Reported    *Reported  `json:"reported"`
	DueDate     *time.Time `json:"dueDate"`
	Assignee    *Assignee  `json:"assignee"`
	Description string     `json:"description"`
	Title       string     `json:"title"`
}

type SearchResponse struct {
	Index int     `json:"index"`
	Total int     `json:"total"`
	Tasks []*Task `json:"tasks"`
}
