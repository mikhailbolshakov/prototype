package grpc

import (
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

	return &domain.Task{
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
	}
}

func (s *Server) fromDomain(task *domain.Task) *pb.Task {

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
		Details:     task.Details,
		ChannelId:   task.ChannelId,
	}
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
