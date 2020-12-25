package repository

import (
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"time"
)

type Task struct {
	kit.BaseDto
	Id string `gorm:"column:id"`
	Num string `gorm:"column:num"`
	Type string `gorm:"column:type"`
	SubType string `gorm:"column:subtype"`
	Status string `gorm:"column:status"`
	SubStatus string `gorm:"column:substatus"`
	ReportedBy string `gorm:"column:reported_by"`
	ReportedAt time.Time `gorm:"column:reported_at"`
	DueDate *time.Time `gorm:"column:due_date"`
	AssigneeGroup string `gorm:"column:assignee_group"`
	AssigneeUser string `gorm:"column:assignee_user"`
	AssigneeAt *time.Time `gorm:"column:assignee_at"`
	Description string `gorm:"column:description"`
	Title string `gorm:"column:title"`
	Details string `gorm:"column:details"`
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
