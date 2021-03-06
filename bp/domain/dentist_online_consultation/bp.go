package dentist_online_consultation

import (
	"context"
	"fmt"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	b "gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/bp/logger"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	pbChat "gitlab.medzdrav.ru/prototype/proto/chat"
	services2 "gitlab.medzdrav.ru/prototype/proto/services"
	taskPb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"time"
)

const (
	TT_CLIENT                = "client"
	TST_DENTIST_CONSULTATION = "dentist-consultation"
)

type bpImpl struct {
	taskService b.TaskService
	userService b.UserService
	chatService b.ChatService
	balance     b.BalanceService
	delivery    b.DeliveryService
	utils       *zeebe.Utils
	bpm.BpBase
}

func NewBp(balanceService b.BalanceService,
	delivery b.DeliveryService,
	taskService b.TaskService,
	userService b.UserService,
	chatService b.ChatService,
	bpm bpm.Engine) b.BusinessProcess {

	bp := &bpImpl{
		delivery:    delivery,
		balance:     balanceService,
		taskService: taskService,
		userService: userService,
		chatService: chatService,
		utils:       zeebe.NewUtils(logger.LF()),
	}
	bp.Engine = bpm

	return bp

}

func (bp *bpImpl) Init() error {

	err := bp.registerBpmHandlers()
	if err != nil {
		return err
	}

	return nil
}

func (bp *bpImpl) SetQueueListeners(ql listener.QueueListener) {
	ql.Add(queue.QUEUE_TYPE_AT_LEAST_ONCE, taskPb.QUEUE_TOPIC_TASK_DUEDATE, bp.dueDateTaskHandler)
	ql.Add(queue.QUEUE_TYPE_AT_LEAST_ONCE, taskPb.QUEUE_TOPIC_TASK_SOLVED_STATUS, bp.solvedTaskHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-expert-online-consultation"
}

func (bp *bpImpl) GetBPMNFileName() string {
	return "dentist_online_consultation.bpmn"
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

func (bp *bpImpl) l() log.CLogger {
	return logger.L().Pr("zeebe").Cmp(bp.GetId())
}

func (bp *bpImpl) dueDateTaskHandler(msg []byte) error {

	task := &taskPb.TaskMessagePayload{}
	ctx, err := queue.Decode(nil, msg, task)
	if err != nil {
		return err
	}

	bp.l().Mth("task-duedate").F(log.FF{"task-id": task.Id}).C(ctx).Dbg().Trc(string(msg))

	_ = bp.SendMessage("msg-consultation-time", task.Id, nil)
	return nil
}

func (bp *bpImpl) solvedTaskHandler(msg []byte) error {

	task := &taskPb.TaskMessagePayload{}
	ctx, err := queue.Decode(nil, msg, task)
	if err != nil {
		return err
	}

	bp.l().Mth("task-solved").F(log.FF{"task-id": task.Id}).C(ctx).Dbg().Trc(string(msg))

	if task.Type.Type == TT_CLIENT && task.Type.SubType == TST_DENTIST_CONSULTATION {

		vars := map[string]interface{}{}
		vars["taskCompleted"] = true
		_ = bp.SendMessage("msg-task-finished", task.Id, vars)

		msg := fmt.Sprintf("???????????????????????? %s ??????????????????", task.Num)
		if err := bp.chatService.Post(ctx, msg, task.ChannelId, "", false); err != nil {
			return err
		}

	}

	return nil
}

func (bp *bpImpl) createTaskHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	deliveryId := variables["deliveryId"].(string)

	dl, err := bp.delivery.GetDelivery(ctx, deliveryId)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	user := bp.userService.Get(ctx, dl.UserId)

	startTime := &dl.StartTime
	expertUserId := dl.Details["expertUserId"].(string)
	consultationTime, err := time.Parse(time.RFC3339, dl.Details["consultationTime"].(string))
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	expert := bp.userService.Get(ctx, expertUserId)

	// check if a channel with this expert already exists
	channels, err := bp.chatService.GetChannelsForUserAndExpert(ctx, user.MMId, expert.MMId)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	var chId string
	if channels != nil && len(channels) > 0 {
		chId = channels[0]
	} else {
		//create a channel
		channelId, err := bp.chatService.CreateClientChannel(ctx, &pbChat.CreateClientChannelRequest{
			ChatUserId:  user.MMId,
			DisplayName: "???????????????????????? ??????????????????????",
			Name:        kit.NewId(),
			Subscribers: []string{expert.MMId},
		})
		if err != nil {
			bp.utils.FailJob(client, job, err)
			return
		}
		chId = channelId
	}

	// create a task
	task, err := bp.taskService.New(ctx, &taskPb.NewTaskRequest{
		Type: &taskPb.Type{
			Type:    "client",
			Subtype: "dentist-consultation",
		},
		Reported: &taskPb.Reported{
			UserId: user.Id,
			At:     grpc.TimeToPbTS(startTime),
		},
		Description: "???????????????????????? ?????????????????? ?????? ?????????????????? ?? ??????????????????????????????",
		Title:       "???????????????????????? ???? ????????????????????????",
		DueDate:     grpc.TimeToPbTS(&consultationTime),
		Assignee: &taskPb.Assignee{
			UserId: expert.Id,
			At:     grpc.TimeToPbTS(startTime),
		},
		ChannelId: chId,
		Reminders: []*taskPb.Reminder{
			{
				BeforeDueDate: &taskPb.BeforeDueDate{
					Unit:  "minutes",
					Value: 1,
				},
			},
		},
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	dl.Details["tasks"] = []string{task.Num}
	dl.Details["channels"] = []string{chId}
	_, err = bp.delivery.UpdateDetails(ctx, dl.Id, dl.Details)
	if err != nil {
		bp.utils.FailJob(client, job, err)
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

	time.Sleep(time.Second * 5)

	if err := bp.chatService.PredefinedPost(ctx, task.ChannelId, user.MMId, "client.new-expert-consultation", true, map[string]interface{}{
		"expert.first-name": expert.ExpertDetails.FirstName,
		"expert.last-name":  expert.ExpertDetails.LastName,
		"due-date":          dueDateStr,
		"expert.url":        expert.ExpertDetails.PhotoUrl,
		"expert.photo-url":  expert.ExpertDetails.PhotoUrl,
	}); err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	if err := bp.chatService.PredefinedPost(ctx, task.ChannelId, expert.MMId, "expert.new-expert-consultation", true, map[string]interface{}{
		"client.first-name": user.ClientDetails.FirstName,
		"client.last-name":  user.ClientDetails.LastName,
		"client.phone":      user.ClientDetails.Phone,
		"client.url":        user.ClientDetails.PhotoUrl,
		"client.med-card":   "https://pmed.moi-service.ru/profile/medcard",
		"due-date":          dueDateStr,
	}); err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	err = bp.utils.CompleteJob(client, job, variables)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) cancelConsultationHandler(client worker.JobClient, job entities.Job) {

	jobKey := job.GetKey()

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	deliveryId := variables["deliveryId"].(string)
	taskId := variables["expertTaskId"].(string)

	dl, err := bp.delivery.GetDelivery(ctx, deliveryId)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	// cancel task
	err = bp.taskService.MakeTransition(ctx, &taskPb.MakeTransitionRequest{
		TaskId:       taskId,
		TransitionId: "5",
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	// cancel service delivery
	_, err = bp.delivery.Cancel(ctx, deliveryId, nil)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	// unlocked service on balance
	_, err = bp.balance.CancelLock(ctx, &services2.ChangeServicesRequest{
		UserId:        dl.UserId,
		ServiceTypeId: dl.ServiceTypeId,
		Quantity:      1,
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	_, err = client.NewCompleteJobCommand().JobKey(jobKey).Send(context.Background())
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) clientFeedbackHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	deliveryId := variables["deliveryId"].(string)

	dl, err := bp.delivery.GetDelivery(ctx, deliveryId)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	user := bp.userService.Get(ctx, dl.UserId)
	expertUserId := dl.Details["expertUserId"].(string)
	expertUser := bp.userService.Get(ctx, expertUserId)

	deliveryTasks := dl.Details["tasks"].([]interface{})

	startTime := time.Now().UTC()
	dueDate := startTime.Add(time.Minute * 3)

	// create a task
	t, err := bp.taskService.New(ctx, &taskPb.NewTaskRequest{
		Type: &taskPb.Type{
			Type:    "client",
			Subtype: "client-feedback",
		},
		Reported:    &taskPb.Reported{UserId: expertUser.Id, At: grpc.TimeToPbTS(&startTime)},
		Description: fmt.Sprintf("???????????? ???????? %s %s, ???????????? ?????????????????? ???????????????? ?????????? ?? ???????????????????????? ?? ?????????????????? %s %s", user.ClientDetails.FirstName, user.ClientDetails.LastName, expertUser.ExpertDetails.FirstName, expertUser.ExpertDetails.LastName),
		Title:       fmt.Sprintf("???????????????? ?????????? ?? ???????????????????????? %s", deliveryTasks[0].(string)),
		Assignee: &taskPb.Assignee{
			UserId: user.Id,
			At:     grpc.TimeToPbTS(&startTime),
		},
		DueDate:   grpc.TimeToPbTS(&dueDate),
		ChannelId: user.ClientDetails.CommonChannelId,
		Reminders: []*taskPb.Reminder{
			{
				BeforeDueDate: &taskPb.BeforeDueDate{
					Unit:  "minutes",
					Value: 1,
				},
			},
		},
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	deliveryTasks = append(deliveryTasks, t.Num)
	dl.Details["tasks"] = deliveryTasks
	_, err = bp.delivery.UpdateDetails(ctx, dl.Id, dl.Details)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	variables["feedbackTaskNum"] = t.Num
	err = bp.utils.CompleteJob(client, job, variables)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) taskInProgressHandler(client worker.JobClient, job entities.Job) {

	jobKey := job.GetKey()

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	taskId := variables["expertTaskId"].(string)

	err = bp.taskService.MakeTransition(ctx, &taskPb.MakeTransitionRequest{
		TaskId:       taskId,
		TransitionId: "2",
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	_, err = client.NewCompleteJobCommand().JobKey(jobKey).Send(ctx)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) deliveryCompletedHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	deliveryId := variables["deliveryId"].(string)

	_, err = bp.delivery.Complete(ctx, deliveryId, nil)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	err = bp.utils.CompleteJob(client, job, nil)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}
