package client_request

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	b "gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	chatPb "gitlab.medzdrav.ru/prototype/proto/chat"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"log"
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
	ql.Add("tasks.assigned", bp.TaskAssignedMessageHandler)
	ql.Add("tasks.solved", bp.TaskSolvedMessageHandler)
	ql.Add("mm.posts", bp.MattermostPostMessageHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-client-request"
}

func (bp *bpImpl) GetBPMNPath() string {
	return "../bp/domain/client_request/bp.bpmn"
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

func (bp *bpImpl) executeBotTaskHandler(client worker.JobClient, job entities.Job) {

	log.Println("executeBotTaskHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
	message := variables["message"].(string)
	channelId := variables["channelId"].(string)

	rs, err := bp.chatService.AskBot(&chatPb.AskBotRequest{
		Message: message,
		From:    "client",
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	// check if there is open task
	variables["botSucceeded"] = rs.Found

	if rs.Found {
		if err := bp.chatService.Post(rs.Answer, channelId, "", false, true); err != nil {
			zeebe.FailJob(client, job, err)
			return
		}
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) checkClientOpenTaskHandler(client worker.JobClient, job entities.Job) {

	log.Println("checkClientOpenTaskHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	channelId := variables["channelId"].(string)
	// retrieves tasks by channel
	ts, err := bp.taskService.Search(&pb.SearchRequest{
		Type: &pb.Type{
			Type:    TASK_TYPE_CLIENT,
			Subtype: TASK_SUBTYPE_COMMON_RQ,
		},
		Status:    &pb.Status{Status: TASK_STATUS_OPEN},
		ChannelId: channelId,
		Paging:    &pb.PagingRequest{Index: 0, Size: 1},
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
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

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createClientRequestTaskHandler(client worker.JobClient, job entities.Job) {

	log.Println("createClientRequestTaskHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	channelId := variables["channelId"].(string)
	userId := variables["userId"].(string)

	postTime := kit.TimeFromMillis(int64(variables["postTime"].(float64)))
	ts, _ := ptypes.TimestampProto(postTime)

	user := bp.userService.Get(userId)

	// create a new task
	createdTask, err := bp.taskService.New(&pb.NewTaskRequest{
		Type: &pb.Type{
			Type:    TASK_TYPE_CLIENT,
			Subtype: TASK_SUBTYPE_COMMON_RQ,
		},
		Reported:    &pb.Reported{UserId: user.Id, At: ts},
		Description: "Клиент обратился в чат",
		Title:       "Обращение клиента по общим вопросам",
		DueDate:     nil,
		Assignee:    &pb.Assignee{},
		ChannelId:   channelId,
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	if err := bp.taskService.MakeTransition(&pb.MakeTransitionRequest{
		TaskId:       createdTask.Id,
		TransitionId: "2",
	}); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	if err := bp.chatService.PredefinedPost(channelId, user.MMId, "client.new-request", true, true, map[string]interface{}{
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

	log.Println("createClientRequestTaskHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	channelId := variables["channelId"].(string)
	assigneeUser := variables["assignee"].(string)

	assignee := bp.userService.Get(assigneeUser)

	if err := bp.chatService.Subscribe(assignee.MMId, channelId); err != nil {
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

	log.Println("sendMessageTaskAssignedHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["userId"].(string)
	assigneeUsername := variables["assignee"].(string)
	channelId := variables["channelId"].(string)
	user := bp.userService.Get(userId)
	assignee := bp.userService.Get(assigneeUsername)

	if err := bp.chatService.PredefinedPost(channelId, user.MMId, "client.request-assigned", true, true, map[string]interface{}{
		"consultant.first-name": assignee.ConsultantDetails.FirstName,
		"consultant.last-name":  assignee.ConsultantDetails.LastName,
		"consultant.url":        assignee.ConsultantDetails.PhotoUrl,
	}); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	if err := bp.chatService.PredefinedPost(channelId, assignee.MMId, "consultant.request-assigned", true, true, map[string]interface{}{
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

	log.Println("sendMessageNoAvailableConsultantHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["userId"].(string)
	channelId := variables["channelId"].(string)
	user := bp.userService.Get(userId)

	if err := bp.chatService.PredefinedPost(channelId, user.MMId, "client.no-consultant-available", true, true, nil); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = zeebe.CompleteJob(client, job, nil)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) MattermostPostMessageHandler(payload []byte) error {

	type postQueueMsg struct {
		Id        string `json:"id"`
		CreateAt  int64  `json:"createAt"`
		UpdateAt  int64  `json:"updateAt"`
		EditAt    int64  `json:"editAt"`
		DeleteAt  int64  `json:"deleteAt"`
		UserId    string `json:"userId"`
		ChannelId string `json:"channelId"`
		Message   string `json:"message"`
		Type      string `json:"type"`
	}

	post := &postQueueMsg{}
	if err := json.Unmarshal(payload, post); err != nil {
		return err
	}

	// get user by MM user id
	user, err := bp.userService.GetByMMId(post.UserId)
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

		_, err := bp.StartProcess("p-client-request", variables)
		if err != nil {
			return err
		}

	}

	return nil

}

func (bp *bpImpl) TaskAssignedMessageHandler(payload []byte) error {

	task := &domain.Task{}
	if err := json.Unmarshal(payload, task); err != nil {
		return err
	}

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_COMMON_RQ && task.Assignee.UserId != "" {
		variables := map[string]interface{}{}
		variables["assignee"] = task.Assignee.UserId
		_ = bp.SendMessage("msg-client-task-assigned", task.Id, variables)
	}

	return nil

}

func (bp *bpImpl) TaskSolvedMessageHandler(payload []byte) error {

	task := &domain.Task{}
	if err := json.Unmarshal(payload, task); err != nil {
		return err
	}

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_COMMON_RQ {

		msg := fmt.Sprintf("Консультация %s завершена", task.Num)
		if err := bp.chatService.Post(msg, task.ChannelId, "", false, true); err != nil {
			log.Println(err)
			return err
		}

		user := bp.userService.Get(task.Reported.UserId)
		if err := bp.chatService.PredefinedPost(task.ChannelId, user.Id, "client.feedback", false, true, nil); err != nil {
			log.Println(err)
			return err
		}

	}

	return nil

}
