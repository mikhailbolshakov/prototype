package client_med_request

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	b "gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	chatPb "gitlab.medzdrav.ru/prototype/proto/chat"
	taskPb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"time"
)

const (
	TASK_TYPE_CLIENT         = "client"
	TASK_SUBTYPE_MED_REQUEST = "med-request"
	TASK_STATUS_OPEN         = "open"
)

type bpImpl struct {
	taskService b.TaskService
	userService b.UserService
	chatService b.ChatService
	bpm.BpBase
}

func NewBp(taskService b.TaskService,
	userService b.UserService,
	chatService b.ChatService,
	bpm bpm.Engine) b.BusinessProcess {

	bp := &bpImpl{
		taskService: taskService,
		userService: userService,
		chatService: chatService,
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
	ql.Add(queue.QUEUE_TYPE_AT_LEAST_ONCE, taskPb.QUEUE_TOPIC_TASK_ASSIGN_STATUS, bp.TaskAssignedMessageHandler)
	ql.Add(queue.QUEUE_TYPE_AT_LEAST_ONCE, taskPb.QUEUE_TOPIC_TASK_SOLVED_STATUS, bp.TaskSolvedMessageHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-client-med-request"
}

func (bp *bpImpl) GetBPMNPath() string {
	return "../bp/domain/client_med_request/bp.bpmn"
}

func (bp *bpImpl) registerBpmHandlers() error {
	return bp.RegisterTaskHandlers(map[string]interface{}{
		"st-check-client-med-channel":        bp.checkClientMedChannelHandler,
		"st-create-client-med-channel":       bp.createClientMedChannelHandler,
		"st-check-client-open-med-task":      bp.checkClientOpenMedTaskHandler,
		"st-create-client-med-req-task":      bp.createClientMedRequestTaskHandler,
		"st-subscribe-med-consultant":        bp.subscribeConsultantHandler,
		"st-msg-med-task-assigned":           bp.sendMessageTaskAssignedHandler,
		"st-msg-no-available-med-consultant": bp.sendMessageNoAvailableConsultantHandler,
	})
}

func (bp *bpImpl) checkClientMedChannelHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	log.L().Pr("zeebe").Cmp(bp.GetId()).Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	userId := variables["userId"].(string)

	user := bp.userService.Get(ctx, userId)
	variables["channelId"] = user.ClientDetails.MedChannelId

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createClientMedChannelHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	log.L().Pr("zeebe").Cmp(bp.GetId()).Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	userId := variables["userId"].(string)
	user := bp.userService.Get(ctx, userId)

	channelId, err := bp.chatService.CreateClientChannel(ctx, &chatPb.CreateClientChannelRequest{
		ChatUserId:  user.MMId,
		DisplayName: "Медицинские консультации",
		Name:        kit.NewId(),
		Subscribers: []string{},
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	time.Sleep(time.Second)

	user.ClientDetails.MedChannelId = channelId
	user, err = bp.userService.SetClientDetails(ctx, user.Id, user.ClientDetails)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	variables["channelId"] = channelId

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) checkClientOpenMedTaskHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	log.L().Pr("zeebe").Cmp(bp.GetId()).Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	channelId := variables["channelId"].(string)
	// retrieves tasks by channel
	ts, err := bp.taskService.Search(ctx, &taskPb.SearchRequest{
		Type: &taskPb.Type{
			Type:    TASK_TYPE_CLIENT,
			Subtype: TASK_SUBTYPE_MED_REQUEST,
		},
		Status:    &taskPb.Status{Status: TASK_STATUS_OPEN},
		ChannelId: channelId,
		Paging:    &taskPb.PagingRequest{Index: 0, Size: 1},
	})

	// check if there is open task
	variables["taskExists"] = len(ts) > 0

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createClientMedRequestTaskHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	log.L().Pr("zeebe").Cmp(bp.GetId()).Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	channelId := variables["channelId"].(string)
	userId := variables["userId"].(string)

	ts, _ := ptypes.TimestampProto(time.Now().UTC())

	user := bp.userService.Get(ctx, userId)

	// create a new task
	createdTask, err := bp.taskService.New(ctx, &taskPb.NewTaskRequest{
		Type: &taskPb.Type{
			Type:    TASK_TYPE_CLIENT,
			Subtype: TASK_SUBTYPE_MED_REQUEST,
		},
		Reported:    &taskPb.Reported{UserId: user.Id, At: ts},
		Description: "Клиент обратился в чат",
		Title:       "Медицинская консультация",
		DueDate:     nil,
		Assignee:    &taskPb.Assignee{},
		ChannelId:   channelId,
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	if err := bp.taskService.MakeTransition(ctx, &taskPb.MakeTransitionRequest{
		TaskId:       createdTask.Id,
		TransitionId: "2",
	}); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	if err := bp.chatService.PredefinedPost(ctx, channelId, user.MMId, "client.new-med-request", true, map[string]interface{}{
		"client.name": fmt.Sprintf("%s", user.ClientDetails.FirstName),
	}); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	variables["taskId"] = createdTask.Id
	variables["taskNum"] = createdTask.Num
	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) subscribeConsultantHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	log.L().Pr("zeebe").Cmp(bp.GetId()).Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	channelId := variables["channelId"].(string)
	assigneeUser := variables["assignee"].(string)

	assignee := bp.userService.Get(ctx, assigneeUser)

	if err := bp.chatService.Subscribe(ctx, assignee.MMId, channelId); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
	time.Sleep(time.Second)

	err = zeebe.CompleteJob(client, job, nil)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) sendMessageTaskAssignedHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	log.L().Pr("zeebe").Cmp(bp.GetId()).Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	userId := variables["userId"].(string)
	assigneeUsername := variables["assignee"].(string)
	channelId := variables["channelId"].(string)
	user := bp.userService.Get(ctx, userId)
	assignee := bp.userService.Get(ctx, assigneeUsername)

	if err := bp.chatService.PredefinedPost(ctx, channelId, user.MMId, "client.request-assigned", true, map[string]interface{}{
		"consultant.first-name": assignee.ConsultantDetails.FirstName,
		"consultant.last-name":  assignee.ConsultantDetails.LastName,
		"consultant.url":        assignee.ConsultantDetails.PhotoUrl,
	}); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	if err := bp.chatService.PredefinedPost(ctx, channelId, assignee.MMId, "consultant.request-assigned", true, map[string]interface{}{
		"client.first-name": user.ClientDetails.FirstName,
		"client.last-name":  user.ClientDetails.LastName,
		"client.phone":      user.ClientDetails.Phone,
		"client.url":        user.ClientDetails.PhotoUrl,
		"client.med-card":   "https://pmed.moi-service.ru/profile/medcard",
	}); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = zeebe.CompleteJob(client, job, nil)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) sendMessageNoAvailableConsultantHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	log.L().Pr("zeebe").Cmp(bp.GetId()).Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	userId := variables["userId"].(string)
	channelId := variables["channelId"].(string)
	user := bp.userService.Get(ctx, userId)

	if err := bp.chatService.PredefinedPost(ctx, channelId, user.MMId, "client.no-consultant-available", true, nil); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = zeebe.CompleteJob(client, job, nil)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) TaskAssignedMessageHandler(msg []byte) error {

	task := &taskPb.TaskMessagePayload{}
	ctx, err := queue.Decode(nil, msg, task)
	if err != nil {
		return err
	}

	log.L().Pr("queue").Cmp(bp.GetId()).Mth("task-assigned").F(log.FF{"task-id": task.Id}).C(ctx).Dbg().Trc(string(msg))

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_MED_REQUEST && task.Assignee.UserId != "" {
		variables := map[string]interface{}{}
		variables["assignee"] = task.Assignee.UserId
		return bp.SendMessage("msg-client-med-task-assigned", task.Id, variables)
	}

	return nil

}

func (bp *bpImpl) TaskSolvedMessageHandler(msg []byte) error {

	task := &taskPb.TaskMessagePayload{}
	ctx, err := queue.Decode(nil, msg, task)
	if err != nil {
		return err
	}

	l := log.L().Pr("queue").Cmp(bp.GetId()).Mth("task-solved").F(log.FF{"task-id": task.Id}).C(ctx).Dbg().Trc(string(msg))

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_MED_REQUEST {

		msg := fmt.Sprintf("Консультация %s завершена", task.Num)
		if err := bp.chatService.Post(ctx, msg, task.ChannelId, "", false); err != nil {
			l.E(err).Err(err)
			return err
		}

	}

	return nil

}
