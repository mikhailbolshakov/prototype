package impl

import (
	"context"
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/common"
	context2 "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"math/rand"
	"time"
)

type serviceImpl struct {
	scheduler    domain.TaskScheduler
	storage      domain.TaskStorage
	config       domain.ConfigService
	usersService domain.UserService
	queue        queue.Queue
	chatService  domain.ChatService
}

func NewTaskService(
	scheduler domain.TaskScheduler,
	storage domain.TaskStorage,
	config domain.ConfigService,
	usersAdapter domain.UserService,
	queue queue.Queue,
	chatService domain.ChatService,
) domain.TaskService {

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

func (t *serviceImpl) remindsSchedulerHandler(ctx context.Context, taskId string) {

	log.Dbg("reminder fired")

	task := t.Get(ctx, taskId)

	if task.ChannelId != "" {

		var msg string
		if task.DueDate != nil {
			duration := task.DueDate.Sub(time.Now().UTC().Round(time.Second))
			msg = fmt.Sprintf("До наступления срока исполнения по задаче %s осталось %v", task.Num, duration)
		} else {
			msg = fmt.Sprintf("Напоминание по задаче %s", task.Num)
		}

		if err := t.chatService.Post(ctx, msg, task.ChannelId, "", false, true); err != nil {
			log.Err(err, true)
			return
		}

	}

}

func (t *serviceImpl) dueDateSchedulerHandler(ctx context.Context, taskId string) {
	log.Dbg("due date fired")

	task := t.Get(ctx, taskId)

	if err := t.queue.Publish(ctx, "tasks.duedate", &queue.Message{Payload: task}); err != nil {
		log.Err(err, true)
		return
	}

	if task.ChannelId != "" {

		dueDateStr := ""
		if task.DueDate != nil {
			dueDateStr = task.DueDate.Format("2006-01-02 15:04:05")
		}

		msg := fmt.Sprintf("Уведомление о наступлении времени решения по задаче %s (%s)", task.Num, dueDateStr)

		if err := t.chatService.Post(ctx, msg, task.ChannelId, "", false, true); err != nil {
			log.Err(err, true)
			return
		}

	}

}

func (t *serviceImpl) New(ctx context.Context, task *domain.Task) (*domain.Task, error) {

	// check configuration by the task type
	cfg, err := t.config.Get(ctx, task.Type)
	if err != nil {
		return nil, err
	}

	tm := time.Now().UTC()

	// get an initial transition
	tr, err := t.config.InitialTransition(ctx, task.Type)
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

	reportedByUser := t.usersService.Get(ctx, task.Reported.UserId, task.Reported.Username)
	if reportedByUser == nil || reportedByUser.Id == "" {
		return nil, fmt.Errorf("cannot find reporter user")
	} else {
		task.Reported.UserId = reportedByUser.Id
		task.Reported.Username = reportedByUser.Username
		task.Reported.Type = reportedByUser.Type
	}

	assigneeUser := t.usersService.Get(ctx, task.Assignee.UserId, task.Assignee.Username)

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
		task.Reminders = []*domain.Reminder{}
	}

	// save to storage
	task, err = t.storage.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	t.putHistory(ctx, task)

	if tr.QueueTopic != "" {
		if err := t.queue.Publish(ctx, tr.QueueTopic, &queue.Message{Payload: task}); err != nil {
			return nil, err
		}
	}

	// add task to scheduler
	if task.DueDate != nil || len(task.Reminders) > 0 {
		t.scheduler.ScheduleTask(ctx, task)
	}

	return task, nil

}

func (t *serviceImpl) putHistory(ctx context.Context, task *domain.Task) {

	go func() {

		r, ok := context2.Request(ctx)
		if !ok {
			log.Err(fmt.Errorf("invalid context"), true)
		}

		hist := &domain.History{
			Id:       kit.NewId(),
			TaskId:   task.Id,
			Status:   task.Status,
			Assignee: task.Assignee,
			ChangedBy: r.GetUsername(),
			ChangedAt: time.Now().UTC(),
		}

		if _, err := t.storage.CreateHistory(ctx, hist); err != nil {
			log.Err(err, true)
		}

	}()
}

func (t *serviceImpl) GetHistory(ctx context.Context, taskId string) []*domain.History {
	task := t.storage.Get(ctx, taskId)
	return t.storage.GetHistory(ctx, task.Id)
}

func (t *serviceImpl) newNum(cfg *domain.Config) (string, error) {
	// TODO:
	return fmt.Sprintf("%s%d", cfg.NumGenRule.Prefix, rand.Intn(99999)), nil
}

func (t *serviceImpl) MakeTransition(ctx context.Context, taskId, transitionId string) (*domain.Task, error) {

	tm := time.Now().UTC()

	// get task from storage
	task := t.storage.Get(ctx, taskId)
	if task == nil {
		return nil, errors.New(fmt.Sprintf("task not found by id %s", taskId))
	}

	trs, err := t.config.NextTransitions(ctx, task.Type, task.Status)
	if err != nil {
		return nil, err
	}

	var targetTr *domain.Transition
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
	task, err = t.storage.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	t.putHistory(ctx, task)

	if targetTr.QueueTopic != "" {
		if err := t.queue.Publish(ctx, targetTr.QueueTopic, &queue.Message{Payload: task}); err != nil {
			return nil, err
		}
	}

	return task, nil

}

func (t *serviceImpl) setAssignee(ctx context.Context, task *domain.Task, assignee *domain.Assignee) error {

	if assignee.UserId != "" || assignee.Username != "" {
		assigneeUser := t.usersService.Get(ctx, assignee.UserId, assignee.Username)
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

func (t *serviceImpl) SetAssignee(ctx context.Context, taskId string, assignee *domain.Assignee) (*domain.Task, error) {

	task := t.Get(ctx, taskId)
	if task == nil {
		return nil, fmt.Errorf("task not found id = %s", taskId)
	}

	if err := t.setAssignee(ctx, task, assignee); err != nil {
		return nil, err
	}

	task, err := t.update(ctx, task)
	if err != nil {
		return nil, err
	}

	t.putHistory(ctx, task)

	return task, nil
}

func (t *serviceImpl) Get(ctx context.Context, taskId string) *domain.Task {
	return t.storage.Get(ctx, taskId)
}

func (t *serviceImpl) update(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	return t.storage.Update(ctx, task)
}

func (t *serviceImpl) Update(ctx context.Context, task *domain.Task) (*domain.Task, error) {

	task, err := t.update(ctx, task)
	if err != nil {
		return nil, err
	}

	t.putHistory(ctx, task)

	return task, nil
}

func (t *serviceImpl) GetByChannel(ctx context.Context, channelId string) []*domain.Task {
	return t.storage.GetByChannel(ctx, channelId)
}

func (t *serviceImpl) GetAssignmentLog(ctx context.Context, cr *domain.AssignmentLogCriteria) (*domain.AssignmentLogResponse, error) {

	if cr.PagingRequest == nil {
		cr.PagingRequest = &common.PagingRequest{}
	}

	if cr.Size == 0 {
		cr.Size = 100
	}

	if cr.Index == 0 {
		cr.Index = 1
	}

	return t.storage.GetAssignmentLog(ctx, cr)
}

func (t *serviceImpl) Search(ctx context.Context, cr *domain.SearchCriteria) (*domain.SearchResponse, error) {

	if cr.PagingRequest == nil {
		cr.PagingRequest = &common.PagingRequest{}
	}

	if cr.Size == 0 {
		cr.Size = 100
	}

	if cr.Index == 0 {
		cr.Index = 1
	}

	return t.storage.Search(ctx, cr)
}

