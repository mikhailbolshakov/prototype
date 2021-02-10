package bp

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/bp/domain/client_law_request"
	"gitlab.medzdrav.ru/prototype/bp/domain/client_med_request"
	"gitlab.medzdrav.ru/prototype/bp/domain/client_request"
	"gitlab.medzdrav.ru/prototype/bp/domain/create_user"
	"gitlab.medzdrav.ru/prototype/bp/domain/dentist_online_consultation"
	"gitlab.medzdrav.ru/prototype/bp/grpc"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/bp_engine"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/config"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/keycloak"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/services"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"math/rand"
)

type serviceImpl struct {
	tasksAdapter    tasks.Adapter
	taskService     domain.TaskService
	usersAdapter    users.Adapter
	usersService    domain.UserService
	servicesAdapter services.Adapter
	chatAdapter     chat.Adapter
	chatService     domain.ChatService
	configAdapter   config.Adapter
	configService   domain.ConfigService
	bps             []domain.BusinessProcess
	queue           queue.Queue
	queueListener   listener.QueueListener
	bpEngineAdapter bp_engine.Adapter
	keycloakAdapter keycloak.Adapter
	grpc            *grpc.Server
}

func New() service.Service {

	s := &serviceImpl{}

	s.bpEngineAdapter = bp_engine.NewAdapter()

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

	s.keycloakAdapter = keycloak.NewAdapter()

	engine := s.bpEngineAdapter.GetEngine()

	// register business processes
	s.bps = append([]domain.BusinessProcess{}, dentist_online_consultation.NewBp(s.servicesAdapter.GetBalanceService(),
		s.servicesAdapter.GetDeliveryService(),
		s.taskService, s.usersService, s.chatService, engine))
	s.bps = append(s.bps, client_request.NewBp(s.taskService, s.usersService, s.chatService, engine))
	s.bps = append(s.bps, create_user.NewBp(s.usersService, s.chatService, engine, s.keycloakAdapter.GetProvider()))
	s.bps = append(s.bps, client_med_request.NewBp(s.taskService, s.usersService, s.chatService, engine))
	s.bps = append(s.bps, client_law_request.NewBp(s.taskService, s.usersService, s.chatService, engine))

	s.grpc = grpc.New(engine)

	return s
}

func (s *serviceImpl) Init() error {

	if err := s.configAdapter.Init(); err != nil {
		return err
	}

	c, err := s.configService.Get()
	if err != nil {
		return err
	}

	if err := s.keycloakAdapter.Init(c); err != nil {
		return err
	}

	if err := s.bpEngineAdapter.Init(c, s.bps, s.queueListener); err != nil {
		return err
	}

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
	s.bpEngineAdapter.Close()
	s.grpc.Close()
	_ = s.queue.Close()
	s.keycloakAdapter.Close()

}
