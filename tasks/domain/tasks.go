package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/chat"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
	"log"
	"math/rand"
	"time"
)

type TaskService interface {
	// register a new task
	New(task *Task) (*Task, error)
	// set status
	MakeTransition(taskId, transitionId string) (*Task, error)
	// assign
	SetAssignee(taskId string, target *Assignee) (*Task, error)
	// get by Id
	Get(taskId string) *Task
	// get tasks by channel
	GetByChannel(channelId string) []*Task
	// update task
	Update(task *Task) (*Task, error)
	// get assignment tasks execution log
	GetAssignmentLog(cr *AssignmentLogCriteria) (*AssignmentLogResponse, error)
	// get task history
	GetHistory(taskId string) []*History
}

type serviceImpl struct {
	scheduler    TaskScheduler
	storage      storage.TaskStorage
	config       ConfigService
	usersService users.Service
	queue        queue.Queue
	chatService  chat.Service
}

func NewTaskService(
	scheduler TaskScheduler,
	storage storage.TaskStorage,
	config ConfigService,
	usersAdapter users.Service,
	queue queue.Queue,
	chatService chat.Service,
) TaskService {

	s := &serviceImpl{
		scheduler:    scheduler,
		storage:      storage,
		config:       config,
		usersService: usersAdapter,
		queue:        queue,
		chatService:  chatService,
	}

	s.scheduler.SetDueDateHandler(s.dueDateSchedulerHandler)
	s.scheduler.SetReminderHandler(s.remindsSchedulerHandler)

	return s
}

func (t *Type) equals(another *Type) bool {
	if another == nil {
		return false
	}

	return t.Type == another.Type && t.SubType == another.SubType
}

func (s *Status) equals(another *Status) bool {
	if another == nil {
		return false
	}

	return s.Status == another.Status && s.SubStatus == another.SubStatus
}

func (t *serviceImpl) remindsSchedulerHandler(taskId string) {

	log.Println("reminder fired")

	task := t.Get(taskId)

	if task.ChannelId != "" {

		var msg string
		if task.DueDate != nil {
			duration := task.DueDate.Sub(time.Now().UTC().Round(time.Second))
			msg = fmt.Sprintf("До наступления срока исполнения по задаче %s осталось %v", task.Num, duration)
		} else {
			msg = fmt.Sprintf("Напоминание по задаче %s", task.Num)
		}

		if err := t.chatService.Post(msg, task.ChannelId, "", false, true); err != nil {
			log.Println("ERROR!!!", err)
			return
		}

	}

}

func (t *serviceImpl) dueDateSchedulerHandler(taskId string) {
	log.Println("due date fired")

	task := t.Get(taskId)

	t.publish(task, "tasks.duedate")

	if task.ChannelId != "" {

		dueDateStr := ""
		if task.DueDate != nil {
			dueDateStr = task.DueDate.Format("2006-01-02 15:04:05")
		}

		msg := fmt.Sprintf("Уведомление о наступлении времени решения по задаче %s (%s)", task.Num, dueDateStr)

		if err := t.chatService.Post(msg, task.ChannelId, "", false, true); err != nil {
			log.Println("ERROR!!!", err)
			return
		}

	}

}

func (t *serviceImpl) New(task *Task) (*Task, error) {

	// check configuration by the task type
	cfg, err := t.config.Get(task.Type)
	if err != nil {
		return nil, err
	}

	tm := time.Now().UTC()

	// get an initial transition
	tr, err := t.config.InitialTransition(task.Type)
	if err != nil {
		return nil, err
	}

	task.Id = kit.NewId()
	task.Num, _ = t.newNum(cfg)
	task.Status = tr.To

	// if reported isn't specified setup current
	if task.Reported.At == nil {
		task.Reported.At = &tm
	}

	reportedByUser := t.usersService.Get(task.Reported.UserId, task.Reported.Username)
	if reportedByUser == nil || reportedByUser.Id == "" {
		return nil, fmt.Errorf("cannot find reporter user")
	} else {
		task.Reported.UserId = reportedByUser.Id
		task.Reported.Username = reportedByUser.Username
		task.Reported.Type = reportedByUser.Type
	}

	assigneeUser := t.usersService.Get(task.Assignee.UserId, task.Assignee.Username)

	if assigneeUser != nil && assigneeUser.Id != "" {
		task.Assignee.Type = assigneeUser.Type
		task.Assignee.Username = assigneeUser.Username
		task.Assignee.UserId = assigneeUser.Id

		// TODO:
		task.Assignee.Group = assigneeUser.Groups[0]

	} else {

		// if assigned user is mandatory for the transition, then throw
		if tr.AssignedUserMandatory {
			return nil, fmt.Errorf("task transition is disallowed due to it's configured as assigned user is manadatory")
		}

		task.Assignee.Type = tr.AutoAssignType
		task.Assignee.Group = tr.AutoAssignGroup

		if task.Assignee.Type == "" {
			return nil, errors.New("no user type specified for the task")
		}

		if task.Assignee.Group == "" {
			return nil, errors.New("no user group specified for the task")
		}

	}

	task.Assignee.At = &tm

	if task.Details == nil {
		task.Details = map[string]interface{}{}
	}

	if task.Reminders == nil {
		task.Reminders = []*Reminder{}
	}

	// save to storage
	dto, err := t.storage.Create(toDto(task))
	if err != nil {
		return nil, err
	}

	task = fromDto(dto)

	t.putHistory(task)

	if tr.QueueTopic != "" {
		t.publish(task, tr.QueueTopic)
	}

	// add task to scheduler
	if task.DueDate != nil || len(task.Reminders) > 0 {
		t.scheduler.ScheduleTask(task)
	}

	return task, nil

}

func (t *serviceImpl) putHistory(task *Task) {

	go func() {
		dto := histToDto(&History{
			Id:       kit.NewId(),
			TaskId:   task.Id,
			Status:   task.Status,
			Assignee: task.Assignee,
			// TODO: current user from session
			ChangedBy: "user",
			ChangedAt: time.Now().UTC(),
		})

		if _, err := t.storage.CreateHistory(dto); err != nil {
			log.Fatal(err)
		}

	}()
}

func (t *serviceImpl) GetHistory(taskId string) []*History {

	dto := t.storage.Get(taskId)

	var res []*History
	for _, dto := range t.storage.GetHistory(dto.Id) {
		res = append(res, histFromDto(dto))
	}
	return res
}

func (t *serviceImpl) newNum(cfg *Config) (string, error) {
	return fmt.Sprintf("%s%d", cfg.NumGenRule.Prefix, rand.Intn(99999)), nil
}

func (t *serviceImpl) MakeTransition(taskId, transitionId string) (*Task, error) {

	tm := time.Now().UTC()

	// get task from storage
	dto := t.storage.Get(taskId)
	if dto == nil {
		return nil, errors.New(fmt.Sprintf("task not found by id %s", taskId))
	}
	task := fromDto(dto)

	trs, err := t.config.NextTransitions(task.Type, task.Status)
	if err != nil {
		return nil, err
	}

	var targetTr *Transition
	for _, tr := range trs {
		if tr.Id == transitionId {
			targetTr = tr
			break
		}
	}
	if targetTr == nil {
		return nil, errors.New(fmt.Sprintf("illegal transition %s", transitionId))
	}

	// set new status
	task.Status = targetTr.To

	// check mandatory assigned user
	if targetTr.AssignedUserMandatory && task.Assignee.UserId == "" {
		return nil, fmt.Errorf("task transition is disallowed due to it's configured as assigned user is manadatory")
	}

	if targetTr.AutoAssignType != "" {
		task.Assignee.At = &tm
		task.Assignee.Type = targetTr.AutoAssignType
	}

	if targetTr.AutoAssignGroup != "" {
		task.Assignee.At = &tm
		task.Assignee.Group = targetTr.AutoAssignGroup
	}

	if task.Assignee.Type == "" {
		return nil, errors.New(fmt.Sprintf("task cannot be assigned on the type %s", task.Assignee.Type))
	}

	if task.Assignee.Group == "" {
		return nil, errors.New(fmt.Sprintf("task cannot be assigned on the group %s", task.Assignee.Group))
	}

	// save to storage
	dto, err = t.storage.Update(toDto(task))
	if err != nil {
		return nil, err
	}

	task = fromDto(dto)

	t.putHistory(task)

	if targetTr.QueueTopic != "" {
		t.publish(task, targetTr.QueueTopic)
	}

	return task, nil

}

func (t *serviceImpl) publish(task *Task, topic string) {
	go func() {

		j, err := json.Marshal(t.taskToQueue(task))
		if err != nil {
			log.Fatal(err)
			return
		}
		err = t.queue.Publish(topic, j)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()
}

func (t *serviceImpl) setAssignee(task *Task, assignee *Assignee) error {

	if assignee.UserId != "" || assignee.Username != "" {
		assigneeUser := t.usersService.Get(assignee.UserId, assignee.Username)
		if assigneeUser == nil || assigneeUser.Id == "" {
			return fmt.Errorf("cannot find asignee")
		}
		task.Assignee.Group = assigneeUser.Type
		task.Assignee.Username = assigneeUser.Username
		task.Assignee.UserId = assigneeUser.Id
		task.Assignee.Type = assigneeUser.Type
	} else {
		// if assignee user isn't passed, then check groups
		// if group passed check if it's allowed in transition
		task.Assignee.Username = ""
		task.Assignee.UserId = ""
		task.Assignee.Group = assignee.Group
		task.Assignee.Type = assignee.Type
		if task.Assignee.Group == "" || task.Assignee.Type == "" {
			return fmt.Errorf("empty assigned group or type")
		}
	}
	tm := time.Now().UTC()
	task.Assignee.At = &tm
	return nil

}

func (t *serviceImpl) SetAssignee(taskId string, assignee *Assignee) (*Task, error) {

	task := t.Get(taskId)
	if task == nil {
		return nil, fmt.Errorf("task not found id = %s", taskId)
	}

	if err := t.setAssignee(task, assignee); err != nil {
		return nil, err
	}

	task, err := t.update(task)
	if err != nil {
		return nil, err
	}

	t.putHistory(task)

	return task, nil
}

func (t *serviceImpl) Get(taskId string) *Task {
	return fromDto(t.storage.Get(taskId))
}

func (t *serviceImpl) update(task *Task) (*Task, error) {

	dto, err := t.storage.Update(toDto(task))
	if err != nil {
		return nil, err
	}

	task = fromDto(dto)

	return task, nil
}

func (t *serviceImpl) Update(task *Task) (*Task, error) {

	task, err := t.update(task)
	if err != nil {
		return nil, err
	}

	t.putHistory(task)

	return task, nil
}

func (t *serviceImpl) GetByChannel(channelId string) []*Task {

	dtos := t.storage.GetByChannel(channelId)
	var res []*Task

	for _, d := range dtos {
		res = append(res, fromDto(d))
	}

	return res

}

func (t *serviceImpl) GetAssignmentLog(cr *AssignmentLogCriteria) (*AssignmentLogResponse, error) {

	if cr.PagingRequest == nil {
		cr.PagingRequest = &common.PagingRequest{}
	}

	if cr.Size == 0 {
		cr.Size = 100
	}

	if cr.Index == 0 {
		cr.Index = 1
	}

	r, err := t.storage.GetAssignmentLog(&storage.AssignmentLogCriteria{
		PagingRequest:   cr.PagingRequest,
		StartTimeAfter:  cr.StartTimeAfter,
		StartTimeBefore: cr.StartTimeBefore,
	})
	if err != nil {
		return nil, err
	}

	return assLogRsFromDto(r), nil
}
