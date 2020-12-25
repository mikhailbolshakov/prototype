package tasks

import (
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
)

func (c *controller) toPb(request *NewTaskRequest) (*pb.NewTaskRequest, error) {

	return &pb.NewTaskRequest{
		Type:        &pb.Type{
			Type:    request.Type.Type,
			Subtype: request.Type.SubType,
		},
		ReportedBy:  request.ReportedBy,
		ReportedAt:  grpc.TimeToPbTS(request.ReportedAt),
		Description: request.Description,
		Title:       request.Title,
		DueDate:     grpc.TimeToPbTS(request.DueDate),
		Assignee:    &pb.Assignee{
			Group: request.Assignee.Group,
			User:  request.Assignee.User,
			At:    grpc.TimeToPbTS(request.Assignee.At),
		},
	}, nil
}

func (s *controller) fromPb(response *pb.Task) (*Task, error) {

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
			ReportedBy:  response.ReportedBy,
			ReportedAt:  grpc.PbTSToTime(response.ReportedAt),
			DueDate:     grpc.PbTSToTime(response.DueDate),
			Assignee:    &Assignee{
				Group: response.Assignee.Group,
				User:  response.Assignee.User,
				At:    grpc.PbTSToTime(response.Assignee.At),
			},
			Description: response.Description,
			Title:       response.Title,
			Details:     response.Details,
		},
		nil
}
