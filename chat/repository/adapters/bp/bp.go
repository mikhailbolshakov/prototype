package bp

import (
	"context"
	"encoding/json"
	pb "gitlab.medzdrav.ru/prototype/proto/bp"
)

type Service interface {
	StartProcess(processId string, vars map[string]interface{}) (string, error)
}

type serviceImpl struct {
	pb.ProcessClient
}

func newImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (s *serviceImpl) StartProcess(processId string, vars map[string]interface{}) (string, error) {
	varsB, _ := json.Marshal(vars)
	rs, err := s.ProcessClient.StartProcess(context.Background(), &pb.StartProcessRequest{
		ProcessId: processId,
		Vars:      varsB,
	})
	if err != nil {
		return "", err
	}
	return rs.Id, nil
}

