package tasks

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
)

func (c *ctrlImpl) toPb(r *NewTaskRequest) *pb.NewTaskRequest {

	var details []byte
	if r.Details != nil {
		details, _ = json.Marshal(r.Details)
	}

	var reminders []*pb.Reminder
	if r.Reminders != nil {

		for _, r := range r.Reminders {
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

	t := &pb.NewTaskRequest{
		Type: &pb.Type{
			Type:    r.Type.Type,
			Subtype: r.Type.SubType,
		},
		Reported: &pb.Reported{
			Type:     r.Reported.Type,
			UserId:   r.Reported.UserId,
			Username: r.Reported.Username,
			At:       grpc.TimeToPbTS(r.Reported.At),
		},
		Description: r.Description,
		Title:       r.Title,
		DueDate:     grpc.TimeToPbTS(r.DueDate),
		Assignee: &pb.Assignee{
			Type:     r.Assignee.Type,
			Group:    r.Assignee.Group,
			UserId:   r.Assignee.UserId,
			Username: r.Assignee.Username,
			At:       grpc.TimeToPbTS(r.Assignee.At),
		},
		Reminders: reminders,
		Details:   details,
	}

	return t
}

func (s *ctrlImpl) fromPb(r *pb.Task) *Task {

	details := map[string]interface{}{}

	if r.Details != nil {
		_ = json.Unmarshal(r.Details, &details)
	}

	t := &Task{
		Id:  r.Id,
		Num: r.Num,
		Type: &Type{
			Type:    r.Type.Type,
			SubType: r.Type.Subtype,
		},
		Status: &Status{
			Status:    r.Status.Status,
			SubStatus: r.Status.Substatus,
		},
		Reported: &Reported{
			Type:     r.Reported.Type,
			UserId:   r.Reported.UserId,
			Username: r.Reported.Username,
			At:       grpc.PbTSToTime(r.Reported.At),
		},
		DueDate: grpc.PbTSToTime(r.DueDate),
		Assignee: &Assignee{
			Type:     r.Assignee.Type,
			Group:    r.Assignee.Group,
			UserId:   r.Assignee.UserId,
			Username: r.Assignee.Username,
			At:       grpc.PbTSToTime(r.Assignee.At),
		},
		Description: r.Description,
		Title:       r.Title,
		Details:     details,
		Reminders:   []*Reminder{},
	}

	for _, r := range r.Reminders {

		dr := &Reminder{}

		if r.BeforeDueDate != nil {
			dr.BeforeDueDate = &BeforeDueDate{
				Unit:  r.BeforeDueDate.Unit,
				Value: uint(r.BeforeDueDate.Value),
			}
		}

		if r.SpecificTime != nil {
			dr.SpecificTime = &SpecificTime{At: grpc.PbTSToTime(r.SpecificTime.At)}
		}

		t.Reminders = append(t.Reminders, dr)
	}

	return t
}

func (s *ctrlImpl) searchRsFromPb(rs *pb.SearchResponse) *SearchResponse {
	r := &SearchResponse{
		Index: int(rs.Paging.Index),
		Total: int(rs.Paging.Total),
		Tasks: []*Task{},
	}

	for _, t := range rs.Tasks {
		r.Tasks = append(r.Tasks, s.fromPb(t))
	}

	return r
}

func (s *ctrlImpl) assLogRsFromPb(rs *pb.AssignmentLogResponse) *AssignmentLogResponse {
	r := &AssignmentLogResponse{
		Index: int(rs.Paging.Index),
		Total: int(rs.Paging.Total),
		Logs:  []*AssignmentLog{},
	}

	for _, t := range rs.Logs {
		r.Logs = append(r.Logs, &AssignmentLog{
			Id:              t.Id,
			StartTime:       grpc.PbTSToTime(t.StartTime),
			FinishTime:      grpc.PbTSToTime(t.FinishTime),
			Status:          t.Status,
			RuleCode:        t.RuleCode,
			RuleDescription: t.RuleDescription,
			UsersInPool:     int(t.UsersInPool),
			TasksToAssign:   int(t.TasksToAssign),
			Assigned:        int(t.Assigned),
			Error:           t.Error,
		})
	}

	return r
}

func (s *ctrlImpl) histFromPb(rs *pb.GetHistoryResponse) []*History {
	var res []*History
	for _, h := range rs.Items {
		res = append(res, &History{
			Status:    &Status{
				Status:    h.Status.Status,
				SubStatus: h.Status.Substatus,
			},
			Assignee:  &Assignee{
				Type:     h.Assignee.Type,
				Group:    h.Assignee.Group,
				UserId:   h.Assignee.UserId,
				Username: h.Assignee.Username,
				At:       grpc.PbTSToTime(h.Assignee.At),
			},
			ChangedBy: h.ChangedBy,
			ChangedAt: *grpc.PbTSToTime(h.ChangedAt),
		})
	}
	return res
}