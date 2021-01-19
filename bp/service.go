package bp

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/bp/bpm"
	"gitlab.medzdrav.ru/prototype/bp/bpm/client_request"
	"gitlab.medzdrav.ru/prototype/bp/bpm/expert_online_consultation"
	"gitlab.medzdrav.ru/prototype/bp/infrastructure"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/services"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/users"
	bpmKit "gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"log"
	"math/rand"
)

type serviceImpl struct {
	tasksAdapter    tasks.Adapter
	usersAdapter    users.Adapter
	servicesAdapter services.Adapter
	mmAdapter       mattermost.Adapter
	bps             []bpm.BusinessProcess
	infr            *infrastructure.Container
	queue           queue.Queue
	queueListener   listener.QueueListener
	bpm             bpmKit.Engine
}

func New() service.Service {

	s := &serviceImpl{}
	s.infr = infrastructure.New()

	s.queue = &stan.Stan{}
	s.queueListener = listener.NewQueueListener(s.queue)

	s.bpm = s.infr.Bpm

	s.servicesAdapter = services.NewAdapter()

	s.tasksAdapter = tasks.NewAdapter(s.queue)
	taskService := s.tasksAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	userService := s.usersAdapter.GetService()

	s.mmAdapter = mattermost.NewAdapter()
	mmService := s.mmAdapter.GetService()

	// register business processes
	s.bps = append([]bpm.BusinessProcess{}, expert_online_consultation.NewBp(s.servicesAdapter.GetBalanceService(),
		s.servicesAdapter.GetDeliveryService(),
		taskService, userService, mmService, s.bpm))

	s.bps = append(s.bps, client_request.NewBp(taskService, userService, mmService, s.bpm))

	return s
}

func (s *serviceImpl) Init() error {

	if err := s.infr.Init(); err != nil {
		return err
	}

	if err := s.tasksAdapter.Init(); err != nil {
		return err
	}

	if err := s.usersAdapter.Init(); err != nil {
		return err
	}

	if err := s.mmAdapter.Init(); err != nil {
		return err
	}

	if err := s.servicesAdapter.Init(); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("client_tasks_%d", rand.Intn(99999))); err != nil {
		return err
	}

	for _, bp := range s.bps {
		if err := bp.Init(); err != nil {
			return err
		}
		bp.SetQueueListeners(s.queueListener)
		log.Printf("business process %s initialized", bp.GetId())
	}

	return nil

}

func (s *serviceImpl) Listen() error {
	return nil
}

func (s *serviceImpl) ListenAsync() error {
	s.queueListener.ListenAsync()
	return nil
}
