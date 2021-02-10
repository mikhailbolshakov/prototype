package bp

import (
	"context"
	"encoding/json"
	pb "gitlab.medzdrav.ru/prototype/proto/bp"
)

type serviceImpl struct {
	pb.ProcessClient
}

func newServiceImpl() *serviceImpl {
	a := &serviceImpl{}
	return a
}

func (s *serviceImpl) StartProcess(processId string, vars map[string]interface{}) (string, error) {

	var varsb []byte

	if vars != nil {
		varsb, _ = json.Marshal(vars)
	}

	rs, err := s.ProcessClient.StartProcess(context.Background(), &pb.StartProcessRequest{
		ProcessId: processId,
		Vars:      varsb,
	})
	if err != nil {
		return "", err
	}

	return rs.Id, nil
}

