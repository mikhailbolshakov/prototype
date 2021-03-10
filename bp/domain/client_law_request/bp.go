package client_law_request

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	b "gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/bp/logger"
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
	TASK_SUBTYPE_LAW_REQUEST = "lawyer-request"
	TASK_STATUS_OPEN         = "open"
)

type bpImpl struct {
	taskService b.TaskService
	userService b.UserService
	chatService b.ChatService
	utils       *zeebe.Utils
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
	ql.Add(queue.QUEUE_TYPE_AT_LEAST_ONCE, taskPb.QUEUE_TOPIC_TASK_ASSIGN_STATUS, bp.TaskAssignedMessageHandler)
	ql.Add(queue.QUEUE_TYPE_AT_LEAST_ONCE, taskPb.QUEUE_TOPIC_TASK_SOLVED_STATUS, bp.TaskSolvedMessageHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-client-law-request"
}

func (bp *bpImpl) GetBPMNFileName() string {
	return "client_law_request.bpmn"
}

func (bp *bpImpl) registerBpmHandlers() error {
	return bp.RegisterTaskHandlers(map[string]interface{}{
		"st-check-client-law-channel":        bp.checkClientLawChannelHandler,
		"st-create-client-law-channel":       bp.createClientLawChannelHandler,
		"st-check-client-open-law-task":      bp.checkClientOpenLawTaskHandler,
		"st-create-client-law-req-task":      bp.createClientLawRequestTaskHandler,
		"st-subscribe-law-consultant":        bp.subscribeConsultantHandler,
		"st-msg-law-task-assigned":           bp.sendMessageTaskAssignedHandler,
		"st-msg-no-available-law-consultant": bp.sendMessageNoAvailableConsultantHandler,
	})
}

func (bp *bpImpl) l() log.CLogger {
	return logger.L().Pr("zeebe").Cmp(bp.GetId())
}

func (bp *bpImpl) checkClientLawChannelHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	userId := variables["userId"].(string)

	user := bp.userService.Get(ctx, userId)
	variables["channelId"] = user.ClientDetails.LawChannelId

	err = bp.utils.CompleteJob(client, job, variables)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createClientLawChannelHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	userId := variables["userId"].(string)
	user := bp.userService.Get(ctx, userId)

	channelId, err := bp.chatService.CreateClientChannel(ctx, &chatPb.CreateClientChannelRequest{
		ChatUserId:  user.MMId,
		DisplayName: "Юридические консультации",
		Name:        kit.NewId(),
		Subscribers: []string{},
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	user.ClientDetails.LawChannelId = channelId
	user, err = bp.userService.SetClientDetails(ctx, user.Id, user.ClientDetails)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	variables["channelId"] = channelId

	err = bp.utils.CompleteJob(client, job, variables)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) checkClientOpenLawTaskHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	channelId := variables["channelId"].(string)
	// retrieves tasks by channel
	ts, err := bp.taskService.Search(ctx, &taskPb.SearchRequest{
		Type: &taskPb.Type{
			Type:    TASK_TYPE_CLIENT,
			Subtype: TASK_SUBTYPE_LAW_REQUEST,
		},
		Status:    &taskPb.Status{Status: TASK_STATUS_OPEN},
		ChannelId: channelId,
		Paging:    &taskPb.PagingRequest{Index: 0, Size: 1},
	})

	// check if there is open task
	variables["taskExists"] = len(ts) > 0

	err = bp.utils.CompleteJob(client, job, variables)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createClientLawRequestTaskHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	channelId := variables["channelId"].(string)
	userId := variables["userId"].(string)

	ts, _ := ptypes.TimestampProto(time.Now().UTC())

	user := bp.userService.Get(ctx, userId)

	// create a new task
	createdTask, err := bp.taskService.New(ctx, &taskPb.NewTaskRequest{
		Type: &taskPb.Type{
			Type:    TASK_TYPE_CLIENT,
			Subtype: TASK_SUBTYPE_LAW_REQUEST,
		},
		Reported:    &taskPb.Reported{UserId: user.Id, At: ts},
		Description: "Клиент обратился в чат",
		Title:       "Юридическая консультация",
		DueDate:     nil,
		Assignee:    &taskPb.Assignee{},
		ChannelId:   channelId,
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	if err := bp.taskService.MakeTransition(ctx, &taskPb.MakeTransitionRequest{
		TaskId:       createdTask.Id,
		TransitionId: "2",
	}); err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	if err := bp.chatService.PredefinedPost(ctx, channelId, user.MMId, "client.new-law-request", true, map[string]interface{}{
		"client.name": fmt.Sprintf("%s", user.ClientDetails.FirstName),
	}); err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	variables["taskId"] = createdTask.Id
	variables["taskNum"] = createdTask.Num
	err = bp.utils.CompleteJob(client, job, variables)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) subscribeConsultantHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	channelId := variables["channelId"].(string)
	assigneeUser := variables["assignee"].(string)

	assignee := bp.userService.Get(ctx, assigneeUser)

	if err := bp.chatService.Subscribe(ctx, assignee.MMId, channelId); err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}
	time.Sleep(time.Second)

	err = bp.utils.CompleteJob(client, job, nil)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) sendMessageTaskAssignedHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

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
		bp.utils.FailJob(client, job, err)
		return
	}

	if err := bp.chatService.PredefinedPost(ctx, channelId, assignee.MMId, "consultant.request-assigned", true, map[string]interface{}{
		"client.first-name": user.ClientDetails.FirstName,
		"client.last-name":  user.ClientDetails.LastName,
		"client.phone":      user.ClientDetails.Phone,
		"client.url":        user.ClientDetails.PhotoUrl,
		"client.med-card":   "https://pmed.moi-service.ru/profile/medcard",
	}); err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	err = bp.utils.CompleteJob(client, job, nil)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) sendMessageNoAvailableConsultantHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	userId := variables["userId"].(string)
	channelId := variables["channelId"].(string)
	user := bp.userService.Get(ctx, userId)

	if err := bp.chatService.PredefinedPost(ctx, channelId, user.MMId, "client.no-consultant-available", true, nil); err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	err = bp.utils.CompleteJob(client, job, nil)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) TaskAssignedMessageHandler(msg []byte) error {

	task := &taskPb.TaskMessagePayload{}
	ctx, err := queue.Decode(nil, msg, task)
	if err != nil {
		return err
	}

	bp.l().Mth("task-assigned").F(log.FF{"task-id": task.Id}).C(ctx).Dbg().Trc(string(msg))

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_LAW_REQUEST && task.Assignee.UserId != "" {
		variables := map[string]interface{}{}
		variables["assignee"] = task.Assignee.UserId
		_ = bp.SendMessage("msg-client-law-task-assigned", task.Id, variables)
	}

	return nil

}

func (bp *bpImpl) TaskSolvedMessageHandler(msg []byte) error {

	task := &taskPb.TaskMessagePayload{}
	ctx, err := queue.Decode(nil, msg, task)
	if err != nil {
		return err
	}

	l := bp.l().Mth("task-solved").F(log.FF{"task-id": task.Id}).C(ctx).Dbg().Trc(string(msg))

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_LAW_REQUEST {

		msg := fmt.Sprintf("Консультация %s завершена", task.Num)
		if err := bp.chatService.Post(ctx, msg, task.ChannelId, "", false); err != nil {
			l.E(err).St().Err()
			return err
		}

	}

	return nil

}
