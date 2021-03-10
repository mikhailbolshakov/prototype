package client_request

import (
	"context"
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
	TASK_TYPE_CLIENT       = "client"
	TASK_SUBTYPE_COMMON_RQ = "common-request"
	TASK_STATUS_OPEN       = "open"
	USER_TYPE_CLIENT       = "client"
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
	ql.Add(queue.QUEUE_TYPE_AT_LEAST_ONCE, chatPb.QUEUE_TOPIC_MATTERMOST_POST_MESSAGE, bp.ChatPostMessageHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-client-request"
}

func (bp *bpImpl) GetBPMNFileName() string {
	return "client_request.bpmn"
}

func (bp *bpImpl) registerBpmHandlers() error {
	return bp.RegisterTaskHandlers(map[string]interface{}{
		"st-bot":                         bp.executeBotTaskHandler,
		"st-check-client-open-task":      bp.checkClientOpenTaskHandler,
		"st-create-client-req-task":      bp.createClientRequestTaskHandler,
		"st-subscribe-consultant":        bp.subscribeConsultantHandler,
		"st-msg-task-assigned":           bp.sendMessageTaskAssignedHandler,
		"st-msg-no-available-consultant": bp.sendMessageNoAvailableConsultantHandler,
	})
}

func (bp *bpImpl) l() log.CLogger {
	return logger.L().Pr("zeebe").Cmp(bp.GetId())
}

func (bp *bpImpl) executeBotTaskHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	message := variables["message"].(string)
	channelId := variables["channelId"].(string)

	rs, err := bp.chatService.AskBot(ctx, &chatPb.AskBotRequest{
		Message: message,
		From:    "client",
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	// check if there is open task
	variables["botSucceeded"] = rs.Found

	if rs.Found {
		if err := bp.chatService.Post(ctx, rs.Answer, channelId, "", false); err != nil {
			bp.utils.FailJob(client, job, err)
			return
		}
	}

	err = bp.utils.CompleteJob(client, job, variables)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) checkClientOpenTaskHandler(client worker.JobClient, job entities.Job) {

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
			Subtype: TASK_SUBTYPE_COMMON_RQ,
		},
		Status:    &taskPb.Status{Status: TASK_STATUS_OPEN},
		ChannelId: channelId,
		Paging:    &taskPb.PagingRequest{Index: 0, Size: 1},
	})
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	// check if there is open task
	taskExists := false
	if len(ts) > 0 {

		for _, t := range ts {
			// TODO: it's simplification
			// a correct check should verify if there are no tasks with close time > post time
			// otherwise this post relates to the closed task and somehow hasn't been delivered in time
			if t.Type.Subtype == TASK_SUBTYPE_COMMON_RQ && t.Status.Status == TASK_STATUS_OPEN {
				variables["taskNum"] = t.Num
				taskExists = true
				break
			}
		}

	}
	variables["taskExists"] = taskExists

	err = bp.utils.CompleteJob(client, job, variables)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createClientRequestTaskHandler(client worker.JobClient, job entities.Job) {

	variables, ctx, err := bp.utils.GetVarsAndCtx(job)
	if err != nil {
		bp.utils.FailJob(client, job, err)
		return
	}

	bp.l().Mth(job.Type).C(ctx).Dbg().Trc(job.String())

	channelId := variables["channelId"].(string)
	userId := variables["userId"].(string)

	postTime := kit.TimeFromMillis(int64(variables["postTime"].(float64)))
	ts, _ := ptypes.TimestampProto(postTime)

	user := bp.userService.Get(ctx, userId)

	// create a new task
	createdTask, err := bp.taskService.New(ctx, &taskPb.NewTaskRequest{
		Type: &taskPb.Type{
			Type:    TASK_TYPE_CLIENT,
			Subtype: TASK_SUBTYPE_COMMON_RQ,
		},
		Reported:    &taskPb.Reported{UserId: user.Id, At: ts},
		Description: "Клиент обратился в чат",
		Title:       "Обращение клиента по общим вопросам",
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

	if err := bp.chatService.PredefinedPost(ctx, channelId, user.MMId, "client.new-request", true, map[string]interface{}{
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

func (bp *bpImpl) ChatPostMessageHandler(msg []byte) error {

	post := &chatPb.MattermostPostMessagePayload{}
	ctx, err := queue.Decode(context.Background(), msg, &post)
	if err != nil {
		return err
	}

	bp.l().Mth("chat-message").C(ctx).Dbg().Trc(string(msg))

	// get user by MM user id
	user, err := bp.userService.GetByMMId(ctx, post.UserId)
	if err != nil {
		return err
	}

	if user != nil && user.Type == USER_TYPE_CLIENT && user.ClientDetails.CommonChannelId == post.ChannelId {

		variables := make(map[string]interface{})
		variables["userId"] = user.Id
		variables["chatUserId"] = post.UserId
		variables["username"] = user.Username
		variables["channelId"] = post.ChannelId
		variables["postTime"] = post.CreateAt
		variables["message"] = post.Message

		if err := bp.utils.CtxToVars(ctx, variables); err != nil {
			return err
		}

		_, err := bp.StartProcess("p-client-request", variables)
		if err != nil {
			return err
		}

	}

	return nil

}

func (bp *bpImpl) TaskAssignedMessageHandler(msg []byte) error {

	task := &taskPb.TaskMessagePayload{}
	ctx, err := queue.Decode(nil, msg, task)
	if err != nil {
		return err
	}

	bp.l().Mth("task-assigned").F(log.FF{"task-id": task.Id}).C(ctx).Dbg().Trc(string(msg))

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_COMMON_RQ && task.Assignee.UserId != "" {
		variables := map[string]interface{}{}
		variables["assignee"] = task.Assignee.UserId
		_ = bp.SendMessage("msg-client-task-assigned", task.Id, variables)
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

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_COMMON_RQ {

		msg := fmt.Sprintf("Консультация %s завершена", task.Num)
		if err := bp.chatService.Post(ctx, msg, task.ChannelId, "", false); err != nil {
			l.E(err).St().Err()
			return err
		}

		user := bp.userService.Get(ctx, task.Reported.UserId)
		if err := bp.chatService.PredefinedPost(ctx, task.ChannelId, user.Id, "client.feedback", false, nil); err != nil {
			l.E(err).St().Err()
			return err
		}

	}

	return nil

}
