package grpc

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
)

func (s *Server) toAssigneeDomain(assignee *pb.Assignee) *domain.Assignee {
	return &domain.Assignee{
		Type:     assignee.Type,
		Group:    assignee.Group,
		UserId:   assignee.UserId,
		Username: assignee.Username,
		At:       grpc.PbTSToTime(assignee.At),
	}
}

func (s *Server) toTaskDomain(r *pb.NewTaskRequest) *domain.Task {

	details := map[string]interface{}{}

	if r.Details != nil {
		_ = json.Unmarshal(r.Details, &details)
	}

	t := &domain.Task{
		Type: &domain.Type{
			Type:    r.Type.Type,
			SubType: r.Type.Subtype,
		},
		Status: &domain.Status{},
		Reported: &domain.Reported{
			Type:     r.Reported.Type,
			UserId:   r.Reported.UserId,
			Username: r.Reported.Username,
			At:       grpc.PbTSToTime(r.Reported.At),
		},
		DueDate: grpc.PbTSToTime(r.DueDate),
		Assignee: &domain.Assignee{
			Type:     r.Assignee.Type,
			Group:    r.Assignee.Group,
			UserId:   r.Assignee.UserId,
			Username: r.Assignee.Username,
			At:       grpc.PbTSToTime(r.Assignee.At),
		},
		Description: r.Description,
		Title:       r.Title,
		ChannelId:   r.ChannelId,
		Reminders:   []*domain.Reminder{},
		Details:     details,
	}

	for _, r := range r.Reminders {

		dr := &domain.Reminder{}

		if r.BeforeDueDate != nil {
			dr.BeforeDueDate = &domain.BeforeDueDate{
				Unit:  domain.TimeUnit(r.BeforeDueDate.Unit),
				Value: uint(r.BeforeDueDate.Value),
			}
		}

		if r.SpecificTime != nil {
			dr.SpecificTime = &domain.SpecificTime{At: grpc.PbTSToTime(r.SpecificTime.At)}
		}

		t.Reminders = append(t.Reminders, dr)
	}

	return t
}

func (s *Server) toTaskPb(task *domain.Task) *pb.Task {

	var details []byte
	if task.Details != nil {
		details, _ = json.Marshal(task.Details)
	}

	var reminders []*pb.Reminder
	if task.Reminders != nil {
		for _, r := range task.Reminders {

			rpb := &pb.Reminder{}

			if r.SpecificTime != nil {
				rpb.SpecificTime = &pb.SpecificTime{At: grpc.TimeToPbTS(r.SpecificTime.At)}
			}

			if r.BeforeDueDate != nil {
				rpb.BeforeDueDate = &pb.BeforeDueDate{
					Unit:  string(r.BeforeDueDate.Unit),
					Value: uint32(r.BeforeDueDate.Value),
				}
			}

			reminders = append(reminders, rpb)

		}
	}

	t := &pb.Task{
		Id:  task.Id,
		Num: task.Num,
		Type: &pb.Type{
			Type:    task.Type.Type,
			Subtype: task.Type.SubType,
		},
		Status: &pb.Status{
			Status:    task.Status.Status,
			Substatus: task.Status.SubStatus,
		},
		Reported: &pb.Reported{
			Type:     task.Reported.Type,
			UserId:   task.Reported.UserId,
			Username: task.Reported.Username,
			At:       grpc.TimeToPbTS(task.Reported.At),
		},
		DueDate: grpc.TimeToPbTS(task.DueDate),
		Assignee: &pb.Assignee{
			Type:     task.Assignee.Type,
			Group:    task.Assignee.Group,
			UserId:   task.Assignee.UserId,
			Username: task.Assignee.Username,
			At:       grpc.TimeToPbTS(task.Assignee.At),
		},
		Description: task.Description,
		Title:       task.Title,
		Details:     details,
		ChannelId:   task.ChannelId,
		Reminders:   reminders,
	}

	return t
}

func (s *Server) toSrchRqDomain(pb *pb.SearchRequest) *domain.SearchCriteria {

	if pb == nil {
		return nil
	}

	return &domain.SearchCriteria{
		PagingRequest: &common.PagingRequest{
			Size:  int(pb.Paging.Size),
			Index: int(pb.Paging.Index),
		},
		Status: &domain.Status{
			Status:    pb.Status.Status,
			SubStatus: pb.Status.Substatus,
		},
		Assignee: &domain.Assignee{
			Type:     pb.Assignee.Type,
			Group:    pb.Assignee.Group,
			UserId:   pb.Assignee.UserId,
			Username: pb.Assignee.Username,
			At:       grpc.PbTSToTime(pb.Assignee.At),
		},
		Type: &domain.Type{
			Type:    pb.Type.Type,
			SubType: pb.Type.Subtype,
		},
		Num:       pb.Num,
		ChannelId: pb.ChannelId,
	}

}

func (s *Server) toSrchRsPb(d *domain.SearchResponse) *pb.SearchResponse {

	rs := &pb.SearchResponse{
		Paging: &pb.PagingResponse{
			Total: int32(d.PagingResponse.Total),
			Index: int32(d.PagingResponse.Index),
		},
		Tasks: []*pb.Task{},
	}

	for _, t := range d.Tasks {
		rs.Tasks = append(rs.Tasks, s.toTaskPb(t))
	}

	return rs
}

func (s *Server) toAssignLogDomain(pb *pb.AssignmentLogRequest) *domain.AssignmentLogCriteria {

	if pb == nil {
		return nil
	}

	return &domain.AssignmentLogCriteria{
		PagingRequest: &common.PagingRequest{
			Size:  int(pb.Paging.Size),
			Index: int(pb.Paging.Index),
		},
		StartTimeAfter:  grpc.PbTSToTime(pb.StartTimeAfter),
		StartTimeBefore: grpc.PbTSToTime(pb.StartTimeBefore),
	}

}

func (s *Server) toAssignLogRsPb(d *domain.AssignmentLogResponse) *pb.AssignmentLogResponse {

	rs := &pb.AssignmentLogResponse{
		Paging: &pb.PagingResponse{
			Total: int32(d.PagingResponse.Total),
			Index: int32(d.PagingResponse.Index),
		},
		Logs: []*pb.AssignmentLog{},
	}

	for _, l := range d.Logs {
		rs.Logs = append(rs.Logs, &pb.AssignmentLog{
			Id:              l.Id,
			StartTime:       grpc.TimeToPbTS(&l.StartTime),
			FinishTime:      grpc.TimeToPbTS(l.FinishTime),
			Status:          l.Status,
			RuleCode:        l.RuleCode,
			RuleDescription: l.RuleDescription,
			UsersInPool:     int32(l.UsersInPool),
			TasksToAssign:   int32(l.TasksToAssign),
			Assigned:        int32(l.Assigned),
			Error:           l.Error,
		})
	}

	return rs
}

func (s *Server) toHistoryPb(src *domain.History) *pb.History {
	return &pb.History{
		Id:     src.Id,
		TaskId: src.TaskId,
		Status: &pb.Status{
			Status:    src.Status.Status,
			Substatus: src.Status.SubStatus,
		},
		Assignee: &pb.Assignee{
			Type:     src.Assignee.Type,
			Group:    src.Assignee.Group,
			UserId:   src.Assignee.UserId,
			Username: src.Assignee.Username,
			At:       grpc.TimeToPbTS(src.Assignee.At),
		},
		ChangedBy: src.ChangedBy,
		ChangedAt: grpc.TimeToPbTS(&src.ChangedAt),
	}
}

func (s *Server) toHistoryDomain(src *pb.History) *domain.History {
	return &domain.History{
		Id:     src.Id,
		TaskId: src.TaskId,
		Status: &domain.Status{
			Status:    src.Status.Status,
			SubStatus: src.Status.Substatus,
		},
		Assignee: &domain.Assignee{
			Type:     src.Assignee.Type,
			Group:    src.Assignee.Group,
			UserId:   src.Assignee.UserId,
			Username: src.Assignee.Username,
			At:       grpc.PbTSToTime(src.Assignee.At),
		},
		ChangedBy: src.ChangedBy,
		ChangedAt: *grpc.PbTSToTime(src.ChangedAt),
	}
}
