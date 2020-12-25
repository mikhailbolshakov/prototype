package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/tasks/repository"
)

func (u *TaskServiceImpl) toDto(domain *Task) (*repository.Task, error) {

	return &repository.Task{
		BaseDto:       storage.BaseDto{},
		Id:            domain.Id,
		Num:           domain.Num,
		Type:          domain.Type.Type,
		SubType:       domain.Type.SubType,
		Status:        domain.Status.Status,
		SubStatus:     domain.Status.SubStatus,
		ReportedBy:    domain.ReportedBy,
		ReportedAt:    *domain.ReportedAt,
		DueDate:       domain.DueDate,
		AssigneeGroup: domain.Assignee.Group,
		AssigneeUser:  domain.Assignee.User,
		AssigneeAt:    domain.Assignee.At,
		Description:   domain.Description,
		Title:         domain.Title,
		Details:       domain.Details,
	}, nil
}

func (u *TaskServiceImpl) fromDto(dto *repository.Task) (*Task, error) {

	return &Task{
		Id:  dto.Id,
		Num: dto.Num,
		Type: &Type{
			Type:    dto.Type,
			SubType: dto.SubType,
		},
		Status: &Status{
			Status:    dto.Status,
			SubStatus: dto.SubStatus,
		},
		ReportedBy: dto.ReportedBy,
		ReportedAt: &dto.ReportedAt,
		DueDate:    dto.DueDate,
		Assignee: &Assignee{
			Group: dto.AssigneeGroup,
			User:  dto.AssigneeUser,
			At:    dto.AssigneeAt,
		},
		Description: dto.Description,
		Title:       dto.Title,
		Details:     dto.Details,
	}, nil
}
