package domain

import (
	"encoding/json"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/queue_model"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
)

func toDto(domain *Task) *storage.Task {

	if domain == nil {
		return nil
	}

	detailsB, _ := json.Marshal(domain.Details)
	remindersB, _ := json.Marshal(domain.Reminders)

	return &storage.Task{
		BaseDto:       kit.BaseDto{},
		Id:            domain.Id,
		Num:           domain.Num,
		Type:          domain.Type.Type,
		SubType:       domain.Type.SubType,
		Status:        domain.Status.Status,
		SubStatus:     domain.Status.SubStatus,
		ReportedBy:    domain.Reported.By,
		ReportedAt:    *domain.Reported.At,
		DueDate:       domain.DueDate,
		AssigneeGroup: domain.Assignee.Group,
		AssigneeUser:  domain.Assignee.User,
		AssigneeAt:    domain.Assignee.At,
		Description:   domain.Description,
		Title:         domain.Title,
		Details:       string(detailsB),
		Reminders:     string(remindersB),
		ChannelId:     domain.ChannelId,
	}
}

func fromDto(dto *storage.Task) *Task {

	if dto == nil {
		return nil
	}

	var reminders []*Reminder
	_ = json.Unmarshal([]byte(dto.Reminders), &reminders)

	var details map[string]interface{}
	_ = json.Unmarshal([]byte(dto.Details), &details)

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
		Reported: &Reported{
			By: dto.ReportedBy,
			At: &dto.ReportedAt,
		},
		DueDate: dto.DueDate,
		Assignee: &Assignee{
			Group: dto.AssigneeGroup,
			User:  dto.AssigneeUser,
			At:    dto.AssigneeAt,
		},
		Description: dto.Description,
		Title:       dto.Title,
		Details:     details,
		Reminders:   reminders,
		ChannelId:   dto.ChannelId,
	}
}

func histToDto(h *History) *storage.History {

	if h == nil {
		return nil
	}

	return &storage.History{
		Id:            h.Id,
		TaskId:        h.TaskId,
		Status:        h.Status.Status,
		SubStatus:     h.Status.SubStatus,
		AssigneeGroup: h.Assignee.Group,
		AssigneeUser:  h.Assignee.User,
		AssigneeAt:    h.Assignee.At,
		ChangedBy:     h.ChangedBy,
		ChangedAt:     h.ChangedAt,
	}
}

func histFromDto(h *storage.History) *History {

	if h == nil {
		return nil
	}

	return &History{
		Id:     h.Id,
		TaskId: h.TaskId,
		Status: &Status{
			Status:    h.Status,
			SubStatus: h.SubStatus,
		},
		Assignee: &Assignee{
			Group: h.AssigneeGroup,
			User:  h.AssigneeUser,
			At:    h.AssigneeAt,
		},
		ChangedBy: h.ChangedBy,
		ChangedAt: h.ChangedAt,
	}
}

func criteriaToDto(c *SearchCriteria) *storage.SearchCriteria {
	if c == nil {
		return nil
	}

	return &storage.SearchCriteria{
		PagingRequest: c.PagingRequest,
		Num:           c.Num,
		Status:        c.Status.Status,
		SubStatus:     c.Status.SubStatus,
		AssigneeGroup: c.Assignee.Group,
		AssigneeUser:  c.Assignee.User,
		Type:          c.Type.Type,
		SubType:       c.Type.SubType,
	}
}

func searchRsFromDto(rs *storage.SearchResponse) *SearchResponse {
	if rs == nil {
		return nil
	}

	r := &SearchResponse{
		PagingResponse: rs.PagingResponse,
		Tasks:          []*Task{},
	}

	for _, t := range rs.Tasks {
		r.Tasks = append(r.Tasks, fromDto(t))
	}

	return r

}

func (ts *serviceImpl) taskToQueue(t *Task) *queue_model.Task {

	res := &queue_model.Task{
		Id:  t.Id,
		Num: t.Num,
		Type: &queue_model.Type{
			Type:    t.Type.Type,
			SubType: t.Type.SubType,
		},
		Status: &queue_model.Status{
			Status:    t.Status.Status,
			SubStatus: t.Status.SubStatus,
		},
		Reported: &queue_model.Reported{
			By: t.Reported.By,
			At: t.Reported.At,
		},
		DueDate: t.DueDate,
		Assignee: &queue_model.Assignee{
			Group: t.Assignee.Group,
			User:  t.Assignee.User,
			At:    t.Assignee.At,
		},
		Description: t.Description,
		Title:       t.Title,
		ChannelId:   t.ChannelId,
	}

	return res
}

func assignmentLogToDto(s *AssignmentLog) *storage.AssignmentLog {
	return &storage.AssignmentLog{
		Id:              s.Id,
		StartTime:       s.StartTime,
		FinishTime:      s.FinishTime,
		Status:          s.Status,
		RuleCode:        s.RuleCode,
		RuleDescription: s.RuleDescription,
		UsersInPool:     s.UsersInPool,
		TasksToAssign:   s.TasksToAssign,
		Assigned:        s.Assigned,
		Error:           s.Error,
	}
}

func assignmentLogFromDto(s *storage.AssignmentLog) *AssignmentLog {
	return &AssignmentLog{
		Id:              s.Id,
		StartTime:       s.StartTime,
		FinishTime:      s.FinishTime,
		Status:          s.Status,
		RuleCode:        s.RuleCode,
		RuleDescription: s.RuleDescription,
		UsersInPool:     s.UsersInPool,
		TasksToAssign:   s.TasksToAssign,
		Assigned:        s.Assigned,
		Error:           s.Error,
	}
}

func assLogRsFromDto(rs *storage.AssignmentLogResponse) *AssignmentLogResponse {
	if rs == nil {
		return nil
	}

	r := &AssignmentLogResponse{
		PagingResponse: rs.PagingResponse,
		Logs:           []*AssignmentLog{},
	}

	for _, t := range rs.Logs {
		r.Logs = append(r.Logs, assignmentLogFromDto(t))
	}

	return r

}
