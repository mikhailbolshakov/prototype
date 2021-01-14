package grpc

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
)

func (s *Server) assigneeFromPb(assignee *pb.Assignee) *domain.Assignee {
	return &domain.Assignee{
		Group: assignee.Group,
		User:  assignee.User,
		At:    grpc.PbTSToTime(assignee.At),
	}
}

func (s *Server) fromPb(request *pb.NewTaskRequest) *domain.Task {

	details := map[string]interface{}{}

	if request.Details != nil {
		_ = json.Unmarshal(request.Details, &details)
	}

	t := &domain.Task{
		Type: &domain.Type{
			Type:    request.Type.Type,
			SubType: request.Type.Subtype,
		},
		Status: &domain.Status{},
		Reported: &domain.Reported{
			By: request.ReportedBy,
			At: grpc.PbTSToTime(request.ReportedAt),
		},
		DueDate: grpc.PbTSToTime(request.DueDate),
		Assignee: &domain.Assignee{
			Group: request.Assignee.Group,
			User:  request.Assignee.User,
			At:    grpc.PbTSToTime(request.Assignee.At),
		},
		Description: request.Description,
		Title:       request.Title,
		ChannelId:   request.ChannelId,
		Reminders:   []*domain.Reminder{},
		Details:     details,
	}

	for _, r := range request.Reminders {

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

func (s *Server) fromDomain(task *domain.Task) *pb.Task {

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
		ReportedBy: task.Reported.By,
		ReportedAt: grpc.TimeToPbTS(task.Reported.At),
		DueDate:    grpc.TimeToPbTS(task.DueDate),
		Assignee: &pb.Assignee{
			Group: task.Assignee.Group,
			User:  task.Assignee.User,
			At:    grpc.TimeToPbTS(task.Assignee.At),
		},
		Description: task.Description,
		Title:       task.Title,
		Details:     details,
		ChannelId:   task.ChannelId,
		Reminders:   reminders,
	}

	return t
}

func (s *Server) searchRqFromPb(pb *pb.SearchRequest) *domain.SearchCriteria {

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
			Group: pb.Assignee.Group,
			User:  pb.Assignee.User,
		},
		Type: &domain.Type{
			Type:    pb.Type.Type,
			SubType: pb.Type.Subtype,
		},
		Num: pb.Num,
	}

}

func (s *Server) searchRsFromDomain(d *domain.SearchResponse) *pb.SearchResponse {

	rs := &pb.SearchResponse{
		Paging: &pb.PagingResponse{
			Total: int32(d.PagingResponse.Total),
			Index: int32(d.PagingResponse.Index),
		},
		Tasks: []*pb.Task{},
	}

	for _, t := range d.Tasks {
		rs.Tasks = append(rs.Tasks, s.fromDomain(t))
	}

	return rs
}

func (s *Server) assLogRqFromPb(pb *pb.AssignmentLogRequest) *domain.AssignmentLogCriteria {

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

func (s *Server) assLogRsFromDomain(d *domain.AssignmentLogResponse) *pb.AssignmentLogResponse {

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
