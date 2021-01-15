package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto/mm"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/queue_model"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
	"log"
	"time"
)

type DeliveryService interface {
	// delivery user service
	RegisterBpmHandlers() error
	Delivery(rq *DeliveryRequest) (*Delivery, error)
	Complete(deliveryId string, finishTime *time.Time) (*Delivery, error)
	Cancel(deliveryId string, cancelTime *time.Time) (*Delivery, error)
}

func NewDeliveryService(balanceService UserBalanceService,
	taskService tasks.Service,
	userService users.Service,
	mmService mattermost.Service,
	storage storage.Storage,
	queue queue.Queue,
	bpm bpm.Engine) DeliveryService {

	d := &deliveryServiceImpl{
		balance:     balanceService,
		taskService: taskService,
		userService: userService,
		mmService:   mmService,
		storage:     storage,
		queue:       queue,
		bpm:         bpm,
	}

	d.taskService.SetTaskDueDateHandler(d.dueDateTaskHandler)

	return d
}

type deliveryServiceImpl struct {
	balance     UserBalanceService
	taskService tasks.Service
	userService users.Service
	mmService   mattermost.Service
	queue       queue.Queue
	storage     storage.Storage
	bpm         bpm.Engine
}

type deliveryHandler func(dl *Delivery) (*Delivery, error)

func (d *deliveryServiceImpl) RegisterBpmHandlers() error {
	return d.bpm.RegisterTaskHandlers(map[string]interface{}{
		"st-create-task": d.createTaskHandler,
	})
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
		StartTime:     time.Now().UTC(),
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
	if st, ok := d.balance.GetTypes()[rq.ServiceTypeId]; ok {
		if st.DeliveryWfId != "" {
			// start WF
			variables := make(map[string]interface{})
			variables["deliveryId"] = delivery.Id
			_, err := d.bpm.StartProcess("p-expert-online-consultation", variables)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(fmt.Sprintf("WF process isn't specified for service type %s", rq.ServiceTypeId))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	// update storage
	dto, err = d.storage.UpdateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}
	delivery = d.deliveryFromDto(dto)

	return delivery, nil
}

func (d *deliveryServiceImpl) Complete(deliveryId string, finishTime *time.Time) (*Delivery, error) {

	delivery := d.deliveryFromDto(d.storage.GetDelivery(deliveryId))
	delivery.Status = "completed"
	if finishTime != nil {
		delivery.FinishTime = finishTime
	} else {
		t := time.Now().UTC()
		delivery.FinishTime = &t
	}

	dto, err := d.storage.UpdateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}

	return d.deliveryFromDto(dto), nil

}

func (d *deliveryServiceImpl) Cancel(deliveryId string, cancelTime *time.Time) (*Delivery, error) {

	delivery := d.deliveryFromDto(d.storage.GetDelivery(deliveryId))
	delivery.Status = "canceled"
	if cancelTime != nil {
		delivery.FinishTime = cancelTime
	} else {
		t := time.Now().UTC()
		delivery.FinishTime = &t
	}

	dto, err := d.storage.UpdateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}

	return d.deliveryFromDto(dto), nil
}

func (d *deliveryServiceImpl) dueDateTaskHandler(ts *queue_model.Task) {
	log.Printf("due date task message received %v", ts)
	_ = d.bpm.SendMessage("msg-consultation-time", ts.Id, nil)
}

func (d *deliveryServiceImpl) createTaskHandler(client worker.JobClient, job entities.Job) {

	log.Println("create task handler executed")

	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		d.failJob(client, job, err)
		return
	}

	deliveryId := variables["deliveryId"].(string)

	dl := d.deliveryFromDto(d.storage.GetDelivery(deliveryId))

	user := d.userService.Get(dl.UserId)

	startTime := &dl.StartTime
	expertUserId := dl.Details["expertUserId"].(string)
	consultationTime, err := time.Parse(time.RFC3339, dl.Details["consultationTime"].(string))
	if err != nil {
		d.failJob(client, job, err)
		return
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
		d.failJob(client, job, err)
		return
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
		d.failJob(client, job, err)
		return
	}

	if err := d.taskService.MakeTransition(&pb.MakeTransitionRequest{
		TaskId:       task.Id,
		TransitionId: "3",
	}); err != nil {
		d.failJob(client, job, err)
		return
	}

	dl.Details["tasks"] = []string{task.Num}
	dl.Details["channels"] = []string{chRs.ChannelId}
	_, err = d.storage.UpdateDelivery(d.deliveryToDto(dl))
	if err != nil {
		d.failJob(client, job, err)
		return
	}

	variables["taskCompleted"] = false
	variables["expertTaskId"] = task.Id

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		d.failJob(client, job, err)
		return
	}

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

}

func (d *deliveryServiceImpl) failJob(client worker.JobClient, job entities.Job, err error) {
	log.Printf("failed to complete job %s error %v", job.GetKey(), err)

	ctx := context.Background()
	_, _ = client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).ErrorMessage(err.Error()).Send(ctx)

}
