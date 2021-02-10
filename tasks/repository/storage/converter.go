package storage

import (
	"encoding/json"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
)

func (s *taskStorageImpl) toTaskDto(domain *domain.Task) *task {

	if domain == nil {
		return nil
	}

	detailsB, _ := json.Marshal(domain.Details)
	remindersB, _ := json.Marshal(domain.Reminders)

	return &task{
		BaseDto:          kit.BaseDto{},
		Id:               domain.Id,
		Num:              domain.Num,
		Type:             domain.Type.Type,
		SubType:          domain.Type.SubType,
		Status:           domain.Status.Status,
		SubStatus:        domain.Status.SubStatus,
		ReportedType:     domain.Reported.Type,
		ReportedUserId:   domain.Reported.UserId,
		ReportedUsername: domain.Reported.Username,
		ReportedAt:       *domain.Reported.At,
		DueDate:          domain.DueDate,
		AssigneeType:     domain.Assignee.Type,
		AssigneeGroup:    domain.Assignee.Group,
		AssigneeUserId:   domain.Assignee.UserId,
		AssigneeUsername: domain.Assignee.Username,
		AssigneeAt:       domain.Assignee.At,
		Description:      domain.Description,
		Title:            domain.Title,
		Details:          string(detailsB),
		Reminders:        string(remindersB),
		ChannelId:        domain.ChannelId,
	}
}

func (s *taskStorageImpl) toTaskIndex(domain *domain.Task) *iTask {
	return &iTask{
		Id:               domain.Id,
		Title:            domain.Title,
		Description:      domain.Description,
		Num:              domain.Num,
		Type:             domain.Type.Type,
		SubType:          domain.Type.SubType,
		Status:           domain.Status.Status,
		SubStatus:        domain.Status.SubStatus,
		AssigneeType:     domain.Assignee.Type,
		AssigneeGroup:    domain.Assignee.Group,
		AssigneeUserId:   domain.Assignee.UserId,
		AssigneeUsername: domain.Assignee.Username,
		ChannelId:        domain.ChannelId,
	}
}

func (s *taskStorageImpl) toTaskDomain(dto *task) *domain.Task {

	if dto == nil || dto.Id == "" {
		return nil
	}

	var reminders []*domain.Reminder
	_ = json.Unmarshal([]byte(dto.Reminders), &reminders)

	var details map[string]interface{}
	_ = json.Unmarshal([]byte(dto.Details), &details)

	return &domain.Task{
		Id:  dto.Id,
		Num: dto.Num,
		Type: &domain.Type{
			Type:    dto.Type,
			SubType: dto.SubType,
		},
		Status: &domain.Status{
			Status:    dto.Status,
			SubStatus: dto.SubStatus,
		},
		Reported: &domain.Reported{
			Type:     dto.ReportedType,
			UserId:   dto.ReportedUserId,
			Username: dto.ReportedUsername,
			At:       &dto.ReportedAt,
		},
		DueDate: dto.DueDate,
		Assignee: &domain.Assignee{
			Type:     dto.AssigneeType,
			Group:    dto.AssigneeGroup,
			UserId:   dto.AssigneeUserId,
			Username: dto.AssigneeUsername,
			At:       dto.AssigneeAt,
		},
		Description: dto.Description,
		Title:       dto.Title,
		Details:     details,
		Reminders:   reminders,
		ChannelId:   dto.ChannelId,
	}
}

func (s *taskStorageImpl) toTasksDomain(dtos []*task) []*domain.Task {

	var res []*domain.Task
	for _, d := range dtos {
		res = append(res, s.toTaskDomain(d))
	}
	return res

}

func (s *taskStorageImpl) toHistoryDto(domain *domain.History) *history {

	if domain == nil {
		return nil
	}

	return &history{
		Id:               domain.Id,
		TaskId:           domain.TaskId,
		Status:           domain.Status.Status,
		SubStatus:        domain.Status.SubStatus,
		AssigneeType:     domain.Assignee.Type,
		AssigneeGroup:    domain.Assignee.Group,
		AssigneeUserId:   domain.Assignee.UserId,
		AssigneeUsername: domain.Assignee.Username,
		AssigneeAt:       domain.Assignee.At,
		ChangedBy:        domain.ChangedBy,
		ChangedAt:        domain.ChangedAt,
	}
}

func (s *taskStorageImpl) toHistoryDomain(h *history) *domain.History {

	if h == nil {
		return nil
	}

	return &domain.History{
		Id:     h.Id,
		TaskId: h.TaskId,
		Status: &domain.Status{
			Status:    h.Status,
			SubStatus: h.SubStatus,
		},
		Assignee: &domain.Assignee{
			Type:     h.AssigneeType,
			Group:    h.AssigneeGroup,
			UserId:   h.AssigneeUserId,
			Username: h.AssigneeUsername,
			At:       h.AssigneeAt,
		},
		ChangedBy: h.ChangedBy,
		ChangedAt: h.ChangedAt,
	}
}

func (s *taskStorageImpl) toHistoriesDomain(dtos []*history) []*domain.History {

	var res []*domain.History
	for _, d := range dtos {
		res = append(res, s.toHistoryDomain(d))
	}
	return res

}

func (s *taskStorageImpl) toAssgnLogDto(domain *domain.AssignmentLog) *assignmentLog {
	return &assignmentLog{
		Id:              domain.Id,
		StartTime:       domain.StartTime,
		FinishTime:      domain.FinishTime,
		Status:          domain.Status,
		RuleCode:        domain.RuleCode,
		RuleDescription: domain.RuleDescription,
		UsersInPool:     domain.UsersInPool,
		TasksToAssign:   domain.TasksToAssign,
		Assigned:        domain.Assigned,
		Error:           domain.Error,
	}
}

func (s *taskStorageImpl) toAssgnLogDomain(dto *assignmentLog) *domain.AssignmentLog {
	return &domain.AssignmentLog{
		Id:              dto.Id,
		StartTime:       dto.StartTime,
		FinishTime:      dto.FinishTime,
		Status:          dto.Status,
		RuleCode:        dto.RuleCode,
		RuleDescription: dto.RuleDescription,
		UsersInPool:     dto.UsersInPool,
		TasksToAssign:   dto.TasksToAssign,
		Assigned:        dto.Assigned,
		Error:           dto.Error,
	}
}

func (s *taskStorageImpl) toAssgnLogsDomain(dtos []*assignmentLog) []*domain.AssignmentLog {

	var res []*domain.AssignmentLog
	for _, d := range dtos {
		res = append(res, s.toAssgnLogDomain(d))
	}
	return res

}

