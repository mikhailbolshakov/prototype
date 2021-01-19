package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/queue"
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
}

func NewTaskService(
	scheduler TaskScheduler,
	storage storage.TaskStorage,
	config ConfigService,
	usersAdapter users.Service,
	queue queue.Queue,
) TaskService {

	s := &serviceImpl{
		scheduler:    scheduler,
		storage:      storage,
		config:       config,
		usersService: usersAdapter,
		queue:        queue,
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
	task := fromDto(t.storage.Get(taskId))
	t.publish(task, "tasks.remind")
}

func (t *serviceImpl) dueDateSchedulerHandler(taskId string) {
	log.Println("due date fired")
	task := fromDto(t.storage.Get(taskId))
	t.publish(task, "tasks.duedate")
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

	reportedByUser := t.usersService.GetByUserName(task.Reported.By)
	if reportedByUser == nil || reportedByUser.Id == "" {
		return nil, fmt.Errorf("cannot find reporter username %s", task.Reported.By)
	}

	if task.Assignee.User != "" {
		assigneeUser := t.usersService.GetByUserName(task.Assignee.User)
		if assigneeUser == nil || assigneeUser.Id == "" {
			return nil, fmt.Errorf("cannot find asignee username %s", task.Assignee.User)
		}
		task.Assignee.Group = assigneeUser.Type
	} else {

		// if assigned user is mandatory for the transition, then throw
		if tr.AssignedUserMandatory {
			return nil, fmt.Errorf("task transition is disallowed due to it's configured as assigned user is manadatory")
		}

		// if assignee user isn't passed, then check groups
		// if group passed check if it's allowed in transition
		if task.Assignee.Group != "" && !tr.checkGroup(task.Assignee.Group) {
			return nil, fmt.Errorf("task cannot be assigned on the group %s", task.Assignee.Group)
		} else {
			// otherwise take auto group if specified in transition
			task.Assignee = &Assignee{Group: tr.AutoAssignGroup}
			if task.Assignee.Group == "" {
				return nil, errors.New("no group specified for the task")
			}
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

func (t *serviceImpl) GetHistory(taskId string) []*History{

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
	if targetTr.AssignedUserMandatory && task.Assignee.User == "" {
		return nil, fmt.Errorf("task transition is disallowed due to it's configured as assigned user is manadatory")
	}

	// check assignee group
	if !targetTr.checkGroup(task.Assignee.Group) {
		task.Assignee.Group = targetTr.AutoAssignGroup
		task.Assignee.At = &tm
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

	if assignee.User != "" {
		assigneeUser := t.usersService.GetByUserName(assignee.User)
		if assigneeUser == nil || assigneeUser.Id == "" {
			return fmt.Errorf("cannot find asignee username %s", assignee.User)
		}
		task.Assignee.Group = assigneeUser.Type
		task.Assignee.User = assigneeUser.Username
	} else {
		// if assignee user isn't passed, then check groups
		// if group passed check if it's allowed in transition
		if assignee.Group != "" {
			task.Assignee.Group = assignee.Group
			task.Assignee.User = ""
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
