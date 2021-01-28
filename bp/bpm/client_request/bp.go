package client_request

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	bpm2 "gitlab.medzdrav.ru/prototype/bp/bpm"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/queue_model"
	"log"
)

type bpImpl struct {
	taskService tasks.Service
	userService users.Service
	mmService   mattermost.Service
	bpm.BpBase
}

func NewBp(taskService tasks.Service,
	userService users.Service,
	mmService mattermost.Service,
	bpm bpm.Engine) bpm2.BusinessProcess {

	bp := &bpImpl{
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

	if err := bp.DeployBPMNs([]string{"../bp/bpm/client_request/bp.bpmn"}); err != nil {
		return err
	}

	return nil
}

func (bp *bpImpl) SetQueueListeners(ql listener.QueueListener) {
	ql.Add("tasks.assigned", bp.TaskAssignedMessageHandler)
	ql.Add("tasks.clientrequest.solved", bp.TaskSolvedMessageHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-client-request"
}

func (bp *bpImpl) registerBpmHandlers() error {
	return bp.RegisterTaskHandlers(map[string]interface{}{
		"st-check-client-open-task":      bp.checkClientOpenTaskHandler,
		"st-create-client-req-task":      bp.createClientRequestTaskHandler,
		"st-subscribe-consultant":        bp.subscribeConsultantHandler,
		"st-msg-task-assigned":           bp.sendMessageTaskAssignedHandler,
		"st-msg-no-available-consultant": bp.sendMessageNoAvailableConsultantHandler,
	})
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
	ts := bp.taskService.GetByChannelId(channelId)

	// check if there is open task
	taskExists := false
	if len(ts) > 0 {

		for _, t := range ts {
			// TODO: it's simplification
			// a correct check should verify if there are no tasks with close time > post time
			// otherwise this post relates to the closed task and somehow hasn't been delivered in time
			if t.Type.Subtype == "medical-request" && t.Status.Status == "open" {
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

	postTime := kit.TimeFromMillis(int64(variables["postTime"].(float64)))
	ts, _ := ptypes.TimestampProto(postTime)

	userId := variables["userId"].(string)
	channelId := variables["channelId"].(string)

	user := bp.userService.Get(userId)

	// create a new task
	createdTask, err := bp.taskService.New(&pb.NewTaskRequest{
		Type: &pb.Type{
			Type:    "client",
			Subtype: "medical-request",
		},
		ReportedBy:  user.Username,
		ReportedAt:  ts,
		Description: "Обращение клиента",
		Title:       "Обращение клиента",
		DueDate:     nil,
		Assignee: &pb.Assignee{
			Group: "consultant",
		},
		ChannelId: channelId,
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

	if err := bp.mmService.SendTriggerPost("client.new-request", user.MMId, user.ClientDetails.MMChannelId, map[string]interface{}{
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

	if err := bp.mmService.Subscribe(assignee.MMId, channelId); err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

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

	if err := bp.mmService.SendTriggerPost("client.request-assigned", user.MMId, channelId, map[string]interface{}{
		"consultant.first-name": assignee.ConsultantDetails.FirstName,
		"consultant.last-name":  assignee.ConsultantDetails.LastName,
		"consultant.url":        "https://prodoctorov.ru/media/photo/tula/doctorimage/589564/432638-589564-ezhikov_l.jpg",
	}); err != nil {
		log.Println(err)
		return
	}

	if err := bp.mmService.SendTriggerPost("consultant.request-assigned", assignee.MMId, channelId, map[string]interface{}{
		"client.first-name": user.ClientDetails.FirstName,
		"client.last-name":  user.ClientDetails.LastName,
		"client.phone":      user.ClientDetails.Phone,
		"client.url":        "https://www.kinonews.ru/insimgs/persimg/persimg3150.jpg",
		"client.med-card":   "https://pmed.moi-service.ru/profile/medcard",
	}); err != nil {
		log.Println(err)
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

	if err := bp.mmService.SendTriggerPost("client.no-consultant-available", user.MMId, channelId, nil); err != nil {
		log.Println(err)
		return
	}

	err = zeebe.CompleteJob(client, job, nil)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) TaskAssignedMessageHandler(payload []byte) error {

	task := &queue_model.Task{}
	if err := json.Unmarshal(payload, task); err != nil {
		return err
	}

	if task.Type.Type == "client" && task.Type.SubType == "medical-request" && task.Assignee.User != "" {
		variables := map[string]interface{}{}
		variables["assignee"] = task.Assignee.User
		_ = bp.SendMessage("msg-client-task-assigned", task.Id, variables)
	}

	return nil

}

func (bp *bpImpl) TaskSolvedMessageHandler(payload []byte) error {

	task := &queue_model.Task{}
	if err := json.Unmarshal(payload, task); err != nil {
		return err
	}

	user := bp.userService.Get(task.Reported.By)

	if err := bp.mmService.SendTriggerPost("client.task-solved", user.MMId, task.ChannelId, map[string]interface{}{
		"task-num": task.Num,
	}); err != nil {
		log.Println(err)
		return err
	}

	return nil

}
