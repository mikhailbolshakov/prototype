package bp

import (
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"gitlab.medzdrav.ru/prototype/bp/bpm"
	"gitlab.medzdrav.ru/prototype/bp/bpm/client_law_request"
	"gitlab.medzdrav.ru/prototype/bp/bpm/client_med_request"
	"gitlab.medzdrav.ru/prototype/bp/bpm/client_request"
	"gitlab.medzdrav.ru/prototype/bp/bpm/create_user"
	"gitlab.medzdrav.ru/prototype/bp/bpm/expert_online_consultation"
	"gitlab.medzdrav.ru/prototype/bp/grpc"
	"gitlab.medzdrav.ru/prototype/bp/infrastructure"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/chat"
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
	chatAdapter     chat.Adapter
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

	s.queue = &stan.Stan{}
	s.queueListener = listener.NewQueueListener(s.queue)

	s.keycloak = gocloak.NewClient("http://localhost:8086")

	s.bpm = s.infr.Bpm

	s.servicesAdapter = services.NewAdapter()

	s.tasksAdapter = tasks.NewAdapter(s.queue)
	taskService := s.tasksAdapter.GetService()

	s.usersAdapter = users.NewAdapter()
	userService := s.usersAdapter.GetService()

	s.chatAdapter = chat.NewAdapter()
	chatService := s.chatAdapter.GetService()

	s.grpc = grpc.New(s.bpm)

	// register business processes
	s.bps = append([]bpm.BusinessProcess{}, expert_online_consultation.NewBp(s.servicesAdapter.GetBalanceService(),
		s.servicesAdapter.GetDeliveryService(),
		taskService, userService, chatService, s.bpm))
	s.bps = append(s.bps, client_request.NewBp(taskService, userService, chatService, s.bpm))
	s.bps = append(s.bps, create_user.NewBp(userService, chatService, s.bpm, s.keycloak))
	s.bps = append(s.bps, client_med_request.NewBp(taskService, userService, chatService, s.bpm))
	s.bps = append(s.bps, client_law_request.NewBp(taskService, userService, chatService, s.bpm))

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
		//go func(){
		//	if err := s.bpm.DeployBPMNs(BPMNs); err != nil {
		//		log.Println("ERROR!!!", "BPMN deploy", err.Error())
		//	}
		//}()

		if err := s.bpm.DeployBPMNs(BPMNs); err != nil {
			log.Println("ERROR!!!", "BPMN deploy", err.Error())
		}
	}

	return nil

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

	if err := s.chatAdapter.Init(); err != nil {
		return err
	}

	if err := s.servicesAdapter.Init(); err != nil {
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

	s.usersAdapter.Close()
	s.chatAdapter.Close()
	s.tasksAdapter.Close()
	s.servicesAdapter.Close()
	s.infr.Close()
	s.grpc.Close()
	_ = s.queue.Close()

}
