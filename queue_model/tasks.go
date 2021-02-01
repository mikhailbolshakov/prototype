package queue_model

import "time"

// TODO: remove this file

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
