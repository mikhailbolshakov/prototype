package domain

import (
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto/mm"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
	"time"
)

type DeliveryService interface {
	// delivery user service
	Delivery(rq *DeliveryRequest) (*Delivery, error)
}

func NewDeliveryService(balanceService UserBalanceService,
	taskService tasks.Service,
	userService users.Service,
	mmService mattermost.Service,
	storage storage.Storage,
	queue queue.Queue) DeliveryService {
	return &deliveryServiceImpl{
		balance:     balanceService,
		taskService: taskService,
		userService: userService,
		mmService:   mmService,
		storage:     storage,
		queue:       queue,
	}
}

type deliveryServiceImpl struct {
	balance     UserBalanceService
	taskService tasks.Service
	userService users.Service
	mmService   mattermost.Service
	queue       queue.Queue
	storage     storage.Storage
}

type deliveryHandler func(dl *Delivery) (*Delivery, error)

func (d *deliveryServiceImpl) getHandlers() map[string]deliveryHandler {
	return map[string]deliveryHandler{
		ST_EXPERT_ONLINE_CONSULTATION: d.onlineConsultationHandler,
	}
}

func (d *deliveryServiceImpl) Delivery(rq *DeliveryRequest) (*Delivery, error) {

	st, ok := d.balance.GetTypes()[rq.ServiceTypeId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	if userBalance, err := d.balance.Get(&GetBalanceRequest{UserId: rq.UserId}); err == nil {

		if b, ok := userBalance.Balance[st]; ok {

			if b.Available == 0 {
				return nil, errors.New(fmt.Sprintf("user %s has no available service of type %s; number of locked services is %d", rq.UserId, rq.ServiceTypeId, b.Locked))
			}

		} else {
			return nil, errors.New(fmt.Sprintf("user %s doesn't have service %s on balance", rq.UserId, rq.ServiceTypeId))
		}

	} else {
		return nil, err
	}

	// lock service
	if _, err := d.balance.Lock(&ModifyBalanceRequest{
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      1,
	}); err != nil {
		return nil, err
	}

	// create a delivery object
	delivery := &Delivery{
		Id:            kit.NewId(),
		UserId:        rq.UserId,
		ServiceTypeId: rq.ServiceTypeId,
		Status:        "in-progress",
		StartTime:     time.Now(),
		FinishTime:    nil,
		Details:       rq.Details,
	}

	// save to storage
	dto, err := d.storage.CreateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}
	delivery = d.deliveryFromDto(dto)

	// execute a handler for corresponding task
	delivery, err = d.getHandlers()[rq.ServiceTypeId](delivery)
	if err != nil {
		return nil, err
	}

	// update storage
	dto, err = d.storage.UpdateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}
	delivery = d.deliveryFromDto(dto)

	return delivery, nil
}

func (d *deliveryServiceImpl) onlineConsultationHandler(dl *Delivery) (*Delivery, error) {

	user := d.userService.Get(dl.UserId)

	startTime := &dl.StartTime
	expertUserId := dl.Details["expertUserId"].(string)
	consultationTime, err := time.Parse(time.RFC3339, dl.Details["consultationTime"].(string))
	if err != nil {
		return nil, err
	}

	expert := d.userService.Get(expertUserId)

	//create a channel
	chRs, err := d.mmService.CreateClientChannel(&mm.CreateClientChannelRequest{
		ClientUserId: user.MMId,
		DisplayName:  "Клиент - эксперт",
		Name:         kit.NewId(),
		Subscribers:  []string{expert.MMId},
	})
	if err != nil {
		return nil, err
	}

	// create a task
	task, err := d.taskService.CreateTask(&pb.NewTaskRequest{
		Type: &pb.Type{
			Type:    "client",
			Subtype: "expert-consultation",
		},
		ReportedBy:  user.Username,
		ReportedAt:  grpc.TimeToPbTS(startTime),
		Description: "Консультация с экспертом",
		Title:       "Консультация с экспертом",
		DueDate:     grpc.TimeToPbTS(&consultationTime),
		Assignee: &pb.Assignee{
			Group: expert.Type,
			User:  expert.Username,
			At:    grpc.TimeToPbTS(startTime),
		},
		ChannelId: chRs.ChannelId,
		Reminders: []*pb.Reminder{
			{
				BeforeDueDate: &pb.BeforeDueDate{
					Unit:  "minutes",
					Value: 1,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if err := d.taskService.MakeTransition(&pb.MakeTransitionRequest{
		TaskId:       task.Id,
		TransitionId: "3",
	}); err != nil {
		return nil, err
	}

	dl.Details["tasks"] = []string{task.Id}
	dl.Details["channel"] = []string{chRs.ChannelId}

	//TODO: we should run WF here

	return dl, nil
}
