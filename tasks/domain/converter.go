package domain

import (
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
)

func toDto(domain *Task) *storage.Task {

	if domain == nil {
		return nil
	}

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
		Details:       domain.Details,
		ChannelId:     domain.ChannelId,
	}
}

func fromDto(dto *storage.Task) *Task {

	if dto == nil {
		return nil
	}

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
		DueDate:    dto.DueDate,
		Assignee: &Assignee{
			Group: dto.AssigneeGroup,
			User:  dto.AssigneeUser,
			At:    dto.AssigneeAt,
		},
		Description: dto.Description,
		Title:       dto.Title,
		Details:     dto.Details,
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
		Tasks: []*Task{},
	}

	for _, t := range rs.Tasks {
		r.Tasks = append(r.Tasks, fromDto(t))
	}

	return r

}
