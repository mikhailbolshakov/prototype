package domain

import (
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/tasks/repository"
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
}

type TaskSearchService interface {
	Search(cr *SearchCriteria) (*SearchResponse, error)
}

type TaskServiceImpl struct {
	storage      repository.TaskStorage
	config       TaskConfigService
	usersAdapter repository.UsersServiceAdapter
}

func NewTaskService() TaskService {
	return &TaskServiceImpl{
		storage:      repository.NewStorage(),
		config:       NewTaskConfigService(),
		usersAdapter: repository.NewUsersServiceAdapter(),
	}
}

func (t *TaskServiceImpl) New(task *Task) (*Task, error) {

	// check configuration by the task type
	cfg, err := t.config.Get(task.Type)
	if err != nil {
		return nil, err
	}

	tm := time.Now()

	// get an initial transition
	tr, err := t.config.InitialTransition(task.Type)
	if err != nil {
		return nil, err
	}

	task.Id = kit.NewId()
	task.Num, _ = t.newNum(cfg)
	task.Status = tr.To

	// if reported at isn't specified setup current
	if task.ReportedAt == nil {
		task.ReportedAt = &tm
	}

	reportedByUser := t.usersAdapter.GetByUserName(task.ReportedBy)
	if reportedByUser == nil || reportedByUser.Id == "" {
		return nil, errors.New(fmt.Sprintf("cannot find reporter username %s", task.ReportedBy))
	}

	if task.Assignee.User != "" {
		assigneeUser := t.usersAdapter.GetByUserName(task.Assignee.User)
		if assigneeUser == nil || assigneeUser.Id == "" {
			return nil, errors.New(fmt.Sprintf("cannot find asignee username %s", task.Assignee.User))
		}
		task.Assignee.Group = assigneeUser.Type
	} else {
		// if assignee user isn't passed, then check groups
		// if group passed check if it's allowed in transition
		if task.Assignee.Group != "" && !tr.checkGroup(task.Assignee.Group) {
			return nil, errors.New(fmt.Sprintf("task cannot be assigned on the group %s", task.Assignee.Group))
		} else {
			// otherwise take auto group if specified in transition
			task.Assignee = &Assignee{Group: tr.AutoAssignGroup}
			if task.Assignee.Group == "" {
				return nil, errors.New("no group specified for the task")
			}
		}

	}

	task.Assignee.At = &tm
	task.Details = "{}"

	// save to storage
	dto, err := t.toDto(task)
	if err != nil {
		return nil, err
	}

	dto, err = t.storage.Create(dto)
	if err != nil {
		return nil, err
	}

	return t.fromDto(dto)

}

func (t *TaskServiceImpl) newNum(cfg *Config) (string, error) {
	return fmt.Sprintf("%s%d", cfg.NumGenRule.Prefix, rand.Intn(99999)), nil
}

func (t *TaskServiceImpl) MakeTransition(taskId, transitionId string) (*Task, error) {

	tm := time.Now()

	// get task from storage
	dto := t.storage.Get(taskId)
	if dto == nil {
		return nil, errors.New(fmt.Sprintf("task not found by id %s", taskId))
	}
	task, err := t.fromDto(dto)
	if err != nil {
		return nil, err
	}

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

	// check assignee group
	if !targetTr.checkGroup(task.Assignee.Group) {
		task.Assignee.Group = targetTr.AutoAssignGroup
		task.Assignee.At = &tm
	}

	if task.Assignee.Group == "" {
		return nil, errors.New(fmt.Sprintf("task cannot be assigned on the group %s", task.Assignee.Group))
	}

	// save to storage
	dto, err = t.toDto(task)
	if err != nil {
		return nil, err
	}

	dto, err = t.storage.Update(dto)
	if err != nil {
		return nil, err
	}

	return t.fromDto(dto)

}

func (t *TaskServiceImpl) SetAssignee(taskId string, target *Assignee) (*Task, error) {
	return nil, nil
}

func (t *TaskServiceImpl) Get(taskId string) *Task {
	return nil
}
