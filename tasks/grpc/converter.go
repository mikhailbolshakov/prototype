package grpc

import (
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
)

func (s *Service) fromPb(request *pb.NewTaskRequest) (*domain.Task, error) {

	return &domain.Task{
		Type: &domain.Type{
			Type:    request.Type.Type,
			SubType: request.Type.Subtype,
		},
		Status:     &domain.Status{},
		ReportedBy: request.ReportedBy,
		ReportedAt: grpc.PbTSToTime(request.ReportedAt),
		DueDate:    grpc.PbTSToTime(request.DueDate),
		Assignee: &domain.Assignee{
			Group: request.Assignee.Group,
			User:  request.Assignee.User,
			At:    grpc.PbTSToTime(request.Assignee.At),
		},
		Description: request.Description,
		Title:       request.Title,
	}, nil
}

func (s *Service) fromDomain(task *domain.Task) (*pb.Task, error) {

	return &pb.Task{
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
			ReportedBy: task.ReportedBy,
			ReportedAt: grpc.TimeToPbTS(task.ReportedAt),
			DueDate:    grpc.TimeToPbTS(task.DueDate),
			Assignee: &pb.Assignee{
				Group: task.Assignee.Group,
				User:  task.Assignee.User,
				At:    grpc.TimeToPbTS(task.Assignee.At),
			},
			Description: task.Description,
			Title:       task.Title,
			Details:     task.Details,
		},
		nil
}
