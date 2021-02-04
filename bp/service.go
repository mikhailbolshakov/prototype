package bp

import (
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"gitlab.medzdrav.ru/prototype/bp/bpm"
	"gitlab.medzdrav.ru/prototype/bp/bpm/client_law_request"
	"gitlab.medzdrav.ru/prototype/bp/bpm/client_med_request"
	"gitlab.medzdrav.ru/prototype/bp/bpm/client_request"
	"gitlab.medzdrav.ru/prototype/bp/bpm/create_user"
	"gitlab.medzdrav.ru/prototype/bp/bpm/dentist_online_consultation"
	"gitlab.medzdrav.ru/prototype/bp/grpc"
	"gitlab.medzdrav.ru/prototype/bp/infrastructure"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/config"
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
	taskService     tasks.Service
	usersAdapter    users.Adapter
	usersService    users.Service
	servicesAdapter services.Adapter
	chatAdapter     chat.Adapter
	chatService     chat.Service
	configAdapter   config.Adapter
	configService   config.Service
	bps             []bpm.BusinessProcess
	infr            *infrastructure.Container
	queue           queue.Queue
	queueListener   listener.QueueListener
	bpm             bpmKit.Engine
	keycloak        gocloak.GoCloak
	grpc            *grpc.Server
}

func New() service.Service {

	s := &serviceImpl{}
	s.infr = infrastructure.New()

	s.configAdapter = config.NewAdapter()
	s.configService = s.configAdapter.GetService()

	s.queue = &stan.Stan{}
	s.queueListener = listener.NewQueueListener(s.queue)

	s.servicesAdapter = services.NewAdapter()

	s.tasksAdapter = tasks.NewAdapter(s.queue)
	s.usersAdapter = users.NewAdapter()
	s.chatAdapter = chat.NewAdapter()
	s.taskService = s.tasksAdapter.GetService()
	s.usersService = s.usersAdapter.GetService()
	s.chatService = s.chatAdapter.GetService()

	return s
}

func (s *serviceImpl) initBPM() error {

	var BPMNs []string
	for _, bp := range s.bps {
		if err := bp.Init(); err != nil {
			return err
		}
		BPMNs = append(BPMNs, bp.GetBPMNPath())
		bp.SetQueueListeners(s.queueListener)
		log.Printf("business process %s initialized", bp.GetId())
	}

	if len(BPMNs) > 0 {
		if err := s.bpm.DeployBPMNs(BPMNs); err != nil {
			log.Println("ERROR!!!", "BPMN deploy", err.Error())
		}
	}

	return nil

}

func (s *serviceImpl) Init() error {

	if err := s.configAdapter.Init(); err != nil {
		return err
	}

	c, err := s.configService.Get()
	if err != nil {
		return err
	}

	s.keycloak = gocloak.NewClient(c.Keycloak.Url)

	if err := s.infr.Init(c); err != nil {
		return err
	}

	s.bpm = s.infr.Bpm

	s.grpc = grpc.New(s.bpm)

	// register business processes
	s.bps = append([]bpm.BusinessProcess{}, dentist_online_consultation.NewBp(s.servicesAdapter.GetBalanceService(),
		s.servicesAdapter.GetDeliveryService(),
		s.taskService, s.usersService, s.chatService, s.bpm))
	s.bps = append(s.bps, client_request.NewBp(s.taskService, s.usersService, s.chatService, s.bpm))
	s.bps = append(s.bps, create_user.NewBp(s.usersService, s.chatService, s.bpm, s.keycloak))
	s.bps = append(s.bps, client_med_request.NewBp(s.taskService, s.usersService, s.chatService, s.bpm))
	s.bps = append(s.bps, client_law_request.NewBp(s.taskService, s.usersService, s.chatService, s.bpm))

	if err := s.grpc.Init(c); err != nil {
		return err
	}

	if err := s.tasksAdapter.Init(c); err != nil {
		return err
	}

	if err := s.usersAdapter.Init(c); err != nil {
		return err
	}

	if err := s.chatAdapter.Init(c); err != nil {
		return err
	}

	if err := s.servicesAdapter.Init(c); err != nil {
		return err
	}

	if err := s.queue.Open(fmt.Sprintf("client_tasks_%d", rand.Intn(99999))); err != nil {
		return err
	}

	if err := s.initBPM(); err != nil {
		return err
	}

	return nil

}

func (s *serviceImpl) ListenAsync() error {

	s.grpc.ListenAsync()
	s.queueListener.ListenAsync()
	return nil
}

func (s *serviceImpl) Close() {
	s.configAdapter.Close()
	s.usersAdapter.Close()
	s.chatAdapter.Close()
	s.tasksAdapter.Close()
	s.servicesAdapter.Close()
	s.infr.Close()
	s.grpc.Close()
	_ = s.queue.Close()

}
