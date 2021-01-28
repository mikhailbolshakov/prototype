package expert_online_consultation

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	bpm2 "gitlab.medzdrav.ru/prototype/bp/bpm"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/services"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/proto/mm"
	services2 "gitlab.medzdrav.ru/prototype/proto/services"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/queue_model"
	"log"
	"time"
)

type bpImpl struct {
	balance     services.BalanceService
	taskService tasks.Service
	userService users.Service
	delivery    services.DeliveryService
	mmService   mattermost.Service
	bpm.BpBase
}

func NewBp(balanceService services.BalanceService,
	delivery services.DeliveryService,
	taskService tasks.Service,
	userService users.Service,
	mmService mattermost.Service,
	bpm bpm.Engine) bpm2.BusinessProcess {

	bp := &bpImpl{
		delivery:    delivery,
		balance:     balanceService,
		taskService: taskService,
		userService: userService,
		mmService:   mmService,
	}
	bp.Engine = bpm

	return bp

}

func (bp *bpImpl) Init() error {

	err := bp.registerBpmHandlers()
	if err != nil {
		return err
	}

	if err := bp.DeployBPMNs([]string{"../bp/bpm/expert_online_consultation/bp.bpmn"}); err != nil {
		return err
	}

	return nil
}

func (bp *bpImpl) SetQueueListeners(ql listener.QueueListener) {
	ql.Add("tasks.duedate", bp.dueDateTaskHandler)
	ql.Add("tasks.solved", bp.solvedTaskHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-expert-online-consultation"
}

func (bp *bpImpl) registerBpmHandlers() error {
	return bp.RegisterTaskHandlers(map[string]interface{}{
		"st-create-task":           bp.createTaskHandler,
		"st-task-in-progress":      bp.taskInProgressHandler,
		"st-complete-consultation": bp.deliveryCompletedHandler,
		"st-client-feedback":       bp.clientFeedbackHandler,
		"st-cancel-consultation":   bp.cancelConsultationHandler,
	})
}

func (bp *bpImpl) dueDateTaskHandler(payload []byte) error {
	ts := &queue_model.Task{}
	if err := json.Unmarshal(payload, ts); err != nil {
		return err
	}

	log.Printf("due date task message received %v", ts)
	_ = bp.SendMessage("msg-consultation-time", ts.Id, nil)
	return nil
}

func (bp *bpImpl) solvedTaskHandler(payload []byte) error {
	ts := &queue_model.Task{}
	if err := json.Unmarshal(payload, ts); err != nil {
		return err
	}
	log.Printf("solved task message received %v", ts)
	vars := map[string]interface{}{}
	vars["taskCompleted"] = true
	_ = bp.SendMessage("msg-task-finished", ts.Id, vars)

	user := bp.userService.Get(ts.Reported.By)

	if err := bp.mmService.SendTriggerPost("client.task-solved", user.MMId, ts.ChannelId, map[string]interface{}{
		"task-num": ts.Num,
	}); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (bp *bpImpl) createTaskHandler(client worker.JobClient, job entities.Job) {

	log.Println("createTaskHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	deliveryId := variables["deliveryId"].(string)

	dl, err := bp.delivery.GetDelivery(deliveryId)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	user := bp.userService.Get(dl.UserId)

	startTime := &dl.StartTime
	expertUserId := dl.Details["expertUserId"].(string)
	consultationTime, err := time.Parse(time.RFC3339, dl.Details["consultationTime"].(string))
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	expert := bp.userService.Get(expertUserId)

	// check if a channel with this expert already exists
	channels, err := bp.mmService.GetChannelsForUserAndExpert(user.MMId, expert.MMId)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	var channelId string
	if channels != nil && len(channels) > 0 {
		channelId = channels[0]
	} else {
		//create a channel
		chRs, err := bp.mmService.CreateClientChannel(&mm.CreateClientChannelRequest{
			ClientUserId: user.MMId,
			DisplayName:  "Клиент - эксперт",
			Name:         kit.NewId(),
			Subscribers:  []string{expert.MMId},
		})
		if err != nil {
			zeebe.FailJob(client, job, err)
			return
		}
		channelId = chRs.ChannelId
		time.Sleep(time.Second * 5)
	}

	// create a task
	task, err := bp.taskService.New(&pb.NewTaskRequest{
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
		ChannelId: channelId,
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
		zeebe.FailJob(client, job, err)
		return
	}

	dl.Details["tasks"] = []string{task.Num}
	dl.Details["channels"] = []string{channelId}
	_, err = bp.delivery.UpdateDetails(dl.Id, dl.Details)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	variables["taskCompleted"] = false
	variables["expertTaskId"] = task.Id
	variables["expertTaskNum"] = task.Num

	dueDate := grpc.PbTSToTime(task.DueDate)
	dueDateStr := ""
	if dueDate != nil {
		dueDateStr = dueDate.Format("2006-01-02 15:04:05")
	}

	if err := bp.mmService.SendTriggerPost("client.new-expert-consultation", user.MMId, task.ChannelId, map[string]interface{}{
		"expert.first-name": expert.ExpertDetails.FirstName,
		"expert.last-name":  expert.ExpertDetails.LastName,
		"due-date":          dueDateStr,
		"expert.url":        "https://prodoctorov.ru/tula/vrach/182956-kasatkin/",
		"expert.photo-url":  "https://prodoctorov.ru/media/photo/tula/doctorimage/182956/323195-182956-kasatkin_l.jpg",
	}); err != nil {
		log.Println(err)
		return
	}

	if err := bp.mmService.SendTriggerPost("expert.new-expert-consultation", expert.MMId, task.ChannelId, map[string]interface{}{
		"client.first-name": user.ClientDetails.FirstName,
		"client.last-name":  user.ClientDetails.LastName,
		"client.phone":      user.ClientDetails.Phone,
		"client.url":        "https://www.kinonews.ru/insimgs/persimg/persimg3150.jpg",
		"client.med-card":   "https://pmed.moi-service.ru/profile/medcard",
		"due-date":          dueDateStr,
	}); err != nil {
		log.Println(err)
		return
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) cancelConsultationHandler(client worker.JobClient, job entities.Job) {

	log.Println("cancelConsultationHandler executed")

	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	deliveryId := variables["deliveryId"].(string)
	taskId := variables["expertTaskId"].(string)

	dl, err := bp.delivery.GetDelivery(deliveryId)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	// cancel task
	err = bp.taskService.MakeTransition(&pb.MakeTransitionRequest{
		TaskId:       taskId,
		TransitionId: "5",
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	// cancel service delivery
	_, err = bp.delivery.Cancel(deliveryId, nil)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	// unlocked service on balance
	_, err = bp.balance.CancelLock(&services2.ChangeServicesRequest{
		UserId:        dl.UserId,
		ServiceTypeId: dl.ServiceTypeId,
		Quantity:      1,
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	_, err = client.NewCompleteJobCommand().JobKey(jobKey).Send(context.Background())
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) clientFeedbackHandler(client worker.JobClient, job entities.Job) {

	log.Println("clientFeedbackHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	deliveryId := variables["deliveryId"].(string)

	dl, err := bp.delivery.GetDelivery(deliveryId)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	user := bp.userService.Get(dl.UserId)
	expertUserId := dl.Details["expertUserId"].(string)
	expertUser := bp.userService.Get(expertUserId)

	deliveryTasks := dl.Details["tasks"].([]interface{})

	startTime := time.Now().UTC()
	dueDate := startTime.Add(time.Minute * 3)

	// create a task
	t, err := bp.taskService.New(&pb.NewTaskRequest{
		Type: &pb.Type{
			Type:    "client",
			Subtype: "client-feedback",
		},
		ReportedBy:  expertUser.Username,
		ReportedAt:  grpc.TimeToPbTS(&startTime),
		Description: fmt.Sprintf("Добрый день %s %s, просим заполнить обратную связь о консультации с экспертом %s %s", user.ClientDetails.FirstName, user.ClientDetails.LastName, expertUser.ExpertDetails.FirstName, expertUser.ExpertDetails.LastName),
		Title:       fmt.Sprintf("Обратная связь о консультации %s", deliveryTasks[0].(string)),
		Assignee: &pb.Assignee{
			Group: user.Type,
			User:  user.Username,
			At:    grpc.TimeToPbTS(&startTime),
		},
		DueDate: grpc.TimeToPbTS(&dueDate),
		ChannelId: user.ClientDetails.MMChannelId,
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
		zeebe.FailJob(client, job, err)
		return
	}

	deliveryTasks = append(deliveryTasks, t.Num)
	dl.Details["tasks"] = deliveryTasks
	_, err = bp.delivery.UpdateDetails(dl.Id, dl.Details)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	variables["feedbackTaskNum"] = t.Num
	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (d *bpImpl) taskInProgressHandler(client worker.JobClient, job entities.Job) {

	log.Println("create task handler executed")

	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	taskId := variables["expertTaskId"].(string)

	err = d.taskService.MakeTransition(&pb.MakeTransitionRequest{
		TaskId:       taskId,
		TransitionId: "2",
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	ctx := context.Background()
	_, err = client.NewCompleteJobCommand().JobKey(jobKey).Send(ctx)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) deliveryCompletedHandler(client worker.JobClient, job entities.Job) {

	log.Println("delivery completed executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	deliveryId := variables["deliveryId"].(string)

	_, err = bp.delivery.Complete(deliveryId, nil)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = zeebe.CompleteJob(client, job, nil)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}
