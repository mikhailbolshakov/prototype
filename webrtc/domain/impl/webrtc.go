package impl

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
)

type webrtcImpl struct {
	common.BaseService
	ion     domain.IonService
	storage domain.WebrtcStorage
}

func NewWebrtcService(ion domain.IonService, storage domain.WebrtcStorage, queue queue.Queue) domain.WebrtcService {

	s := &webrtcImpl{
		ion: ion,
		storage: storage,
	}
	s.BaseService = common.BaseService{Queue: queue}

	return s
}
