package client_law_request

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
	TASK_TYPE_CLIENT         = "client"
	TASK_SUBTYPE_LAW_REQUEST = "lawyer-request"
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
	ql.Add("tasks.assigned", bp.TaskAssignedMessageHandler)
	ql.Add("tasks.solved", bp.TaskSolvedMessageHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-client-law-request"
}

func (bp *bpImpl) GetBPMNPath() string {
	return "../bp/domain/client_law_request/bp.bpmn"
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

func (bp *bpImpl) checkClientLawChannelHandler(client worker.JobClient, job entities.Job) {

	log.Println("checkClientLawChannelHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["userId"].(string)

	user := bp.userService.Get(userId)
	variables["channelId"] = user.ClientDetails.LawChannelId

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createClientLawChannelHandler(client worker.JobClient, job entities.Job) {

	log.Println("createClientLawChannelHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["userId"].(string)
	user := bp.userService.Get(userId)

	channelId, err := bp.chatService.CreateClientChannel(&chatPb.CreateClientChannelRequest{
		ClientUserId: user.MMId,
		DisplayName:  "Юридические консультации",
		Name:         kit.NewId(),
		Subscribers:  []string{},
	})
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	user.ClientDetails.LawChannelId = channelId
	user, err = bp.userService.SetClientDetails(user.Id, user.ClientDetails)
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

func (bp *bpImpl) checkClientOpenLawTaskHandler(client worker.JobClient, job entities.Job) {

	log.Println("checkClientOpenLawTaskHandler executed")

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
			Subtype: TASK_SUBTYPE_LAW_REQUEST,
		},
		Status:    &pb.Status{Status: TASK_STATUS_OPEN},
		ChannelId: channelId,
		Paging:    &pb.PagingRequest{Index: 0, Size: 1},
	})

	// check if there is open task
	variables["taskExists"] = len(ts) > 0

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createClientLawRequestTaskHandler(client worker.JobClient, job entities.Job) {

	log.Println("createClientLawRequestTaskHandler executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	channelId := variables["channelId"].(string)
	userId := variables["userId"].(string)

	ts, _ := ptypes.TimestampProto(time.Now().UTC())

	user := bp.userService.Get(userId)

	// create a new task
	createdTask, err := bp.taskService.New(&pb.NewTaskRequest{
		Type: &pb.Type{
			Type:    TASK_TYPE_CLIENT,
			Subtype: TASK_SUBTYPE_LAW_REQUEST,
		},
		Reported:    &pb.Reported{UserId: user.Id, At: ts},
		Description: "Клиент обратился в чат",
		Title:       "Юридическая консультация",
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

	if err := bp.chatService.PredefinedPost(channelId, user.MMId, "client.new-law-request", true, true, map[string]interface{}{
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

func (bp *bpImpl) TaskAssignedMessageHandler(payload []byte) error {

	task := &domain.Task{}
	if err := json.Unmarshal(payload, task); err != nil {
		return err
	}

	log.Println("task %s assigned (client_law_request handler)")

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_LAW_REQUEST && task.Assignee.UserId != "" {
		variables := map[string]interface{}{}
		variables["assignee"] = task.Assignee.UserId
		_ = bp.SendMessage("msg-client-law-task-assigned", task.Id, variables)
	}

	return nil

}

func (bp *bpImpl) TaskSolvedMessageHandler(payload []byte) error {

	task := &domain.Task{}
	if err := json.Unmarshal(payload, task); err != nil {
		return err
	}

	if task.Type.Type == TASK_TYPE_CLIENT && task.Type.SubType == TASK_SUBTYPE_LAW_REQUEST {

		msg := fmt.Sprintf("Консультация %s завершена", task.Num)
		if err := bp.chatService.Post(msg, task.ChannelId, "", false, true); err != nil {
			log.Println(err)
			return err
		}

	}

	return nil

}
