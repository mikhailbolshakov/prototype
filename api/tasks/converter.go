package tasks

import (
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
)

func (c *controller) toPb(request *NewTaskRequest) *pb.NewTaskRequest {

	return &pb.NewTaskRequest{
		Type:        &pb.Type{
			Type:    request.Type.Type,
			Subtype: request.Type.SubType,
		},
		ReportedBy:  request.Reported.By,
		ReportedAt:  grpc.TimeToPbTS(request.Reported.At),
		Description: request.Description,
		Title:       request.Title,
		DueDate:     grpc.TimeToPbTS(request.DueDate),
		Assignee:    &pb.Assignee{
			Group: request.Assignee.Group,
			User:  request.Assignee.User,
			At:    grpc.TimeToPbTS(request.Assignee.At),
		},
	}
}

func (s *controller) fromPb(response *pb.Task) *Task {

	return &Task{
			Id:          response.Id,
			Num:         response.Num,
			Type:        &Type{
				Type:    response.Type.Type,
				SubType: response.Type.Subtype,
			},
			Status:      &Status{
				Status:    response.Status.Status,
				SubStatus: response.Status.Substatus,
			},
			Reported: &Reported{
				By: response.ReportedBy,
				At: grpc.PbTSToTime(response.ReportedAt),
			},
			DueDate:     grpc.PbTSToTime(response.DueDate),
			Assignee:    &Assignee{
				Group: response.Assignee.Group,
				User:  response.Assignee.User,
				At:    grpc.PbTSToTime(response.Assignee.At),
			},
			Description: response.Description,
			Title:       response.Title,
			Details:     response.Details,
		}
}

func (s *controller) searchRsFromPb(rs *pb.SearchResponse) *SearchResponse {
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

func (s *controller) assLogRsFromPb(rs *pb.AssignmentLogResponse) *AssignmentLogResponse {
	r := &AssignmentLogResponse{
		Index: int(rs.Paging.Index),
		Total: int(rs.Paging.Total),
		Logs: []*AssignmentLog{},
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
