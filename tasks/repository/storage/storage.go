package storage

import (
	"context"
	"github.com/google/uuid"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/common"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"math"
	"time"
)

type task struct {
	kitStorage.BaseDto
	Id               string     `gorm:"column:id"`
	Num              string     `gorm:"column:num"`
	Type             string     `gorm:"column:type"`
	SubType          string     `gorm:"column:subtype"`
	Status           string     `gorm:"column:status"`
	SubStatus        string     `gorm:"column:substatus"`
	ReportedType     string     `gorm:"column:reported_type"`
	ReportedUserId   string     `gorm:"column:reported_user_id"`
	ReportedUsername string     `gorm:"column:reported_username"`
	ReportedAt       time.Time  `gorm:"column:reported_at"`
	DueDate          *time.Time `gorm:"column:due_date"`
	AssigneeType     string     `gorm:"column:assignee_type"`
	AssigneeGroup    string     `gorm:"column:assignee_group"`
	AssigneeUserId   string     `gorm:"column:assignee_user_id"`
	AssigneeUsername string     `gorm:"column:assignee_username"`
	AssigneeAt       *time.Time `gorm:"column:assignee_at"`
	Description      string     `gorm:"column:description"`
	Title            string     `gorm:"column:title"`
	Details          string     `gorm:"column:details"`
	Reminders        string     `gorm:"column:reminders"`
	ChannelId        string     `gorm:"column:channel_id"`
}

type history struct {
	Id               string     `gorm:"column:id"`
	TaskId           string     `gorm:"column:task_id"`
	Status           string     `gorm:"column:status"`
	SubStatus        string     `gorm:"column:substatus"`
	AssigneeType     string     `gorm:"column:assignee_type"`
	AssigneeGroup    string     `gorm:"column:assignee_group"`
	AssigneeUserId   string     `gorm:"column:assignee_user_id"`
	AssigneeUsername string     `gorm:"column:assignee_username"`
	AssigneeAt       *time.Time `gorm:"column:assignee_at"`
	ChangedBy        string     `gorm:"column:changed_by"`
	ChangedAt        time.Time  `gorm:"column:changed_at"`
}

type assignmentLog struct {
	Id              string     `gorm:"column:id"`
	StartTime       time.Time  `gorm:"column:start_time"`
	FinishTime      *time.Time `gorm:"column:finish_time"`
	Status          string     `gorm:"column:status"`
	RuleCode        string     `gorm:"column:rule_code"`
	RuleDescription string     `gorm:"column:rule_description"`
	UsersInPool     int        `gorm:"column:users_in_pool"`
	TasksToAssign   int        `gorm:"column:tasks_to_assign"`
	Assigned        int        `gorm:"column:assigned"`
	Error           string     `gorm:"column:error"`
}

type taskStorageImpl struct {
	c *container
}

func newStorage(c *container) *taskStorageImpl {
	return &taskStorageImpl{c: c}
}

func (s *taskStorageImpl) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {

	dto := s.toTaskDto(task)

	t := time.Now().UTC()
	dto.CreatedAt, dto.UpdatedAt = t, t

	// to DB
	result := s.c.Db.Instance.Create(dto)

	if result.Error != nil {
		return nil, result.Error
	}

	// to index
	s.c.Search.IndexAsync("tasks", task.Id, s.toTaskIndex(task))

	return task, nil
}

func (s *taskStorageImpl) Get(ctx context.Context, id string) *domain.Task {

	// TODO: get from Redis

	task := &task{}
	if _, err := uuid.Parse(id); err == nil {
		task.Id = id
		s.c.Db.Instance.First(task)
	} else {
		s.c.Db.Instance.Where("num = ?", id).First(task)
	}

	return s.toTaskDomain(task)
}

func (s *taskStorageImpl) Update(ctx context.Context, task *domain.Task) (*domain.Task, error) {

	dto := s.toTaskDto(task)

	dto.UpdatedAt = time.Now().UTC()

	// to DB
	result := s.c.Db.Instance.Save(dto)

	if result.Error != nil {
		return nil, result.Error
	}

	// to index
	s.c.Search.IndexAsync("tasks", task.Id, s.toTaskIndex(task))

	return task, nil
}

func (s *taskStorageImpl) GetByChannel(ctx context.Context, channelId string) []*domain.Task {
	var tasks []*task
	s.c.Db.Instance.Where("channel_id = ?", channelId).Find(&tasks)
	return s.toTasksDomain(tasks)
}

func (s *taskStorageImpl) GetByIds(ctx context.Context, ids []string) []*domain.Task {
	var tasks []*task
	s.c.Db.Instance.Find(&tasks, ids)
	return s.toTasksDomain(tasks)
}

func (s *taskStorageImpl) CreateHistory(ctx context.Context, h *domain.History) (*domain.History, error) {
	dto := s.toHistoryDto(h)
	result := s.c.Db.Instance.Create(dto)
	if result.Error != nil {
		return nil, result.Error
	}
	return h, nil
}

func (s *taskStorageImpl) GetHistory(ctx context.Context, taskId string) []*domain.History {
	var dtos []*history
	s.c.Db.Instance.Where("task_id = ?", taskId).Order("changed_at desc").Find(&dtos)
	return s.toHistoriesDomain(dtos)
}

func (s *taskStorageImpl) SaveAssignmentLog(ctx context.Context, l *domain.AssignmentLog) (*domain.AssignmentLog, error) {

	dto := s.toAssgnLogDto(l)
	if l.Id == "" {
		id := kit.NewId()
		dto.Id, l.Id = id, id
		s.c.Db.Instance.Create(dto)
	} else {
		s.c.Db.Instance.Save(dto)
	}

	return l, nil
}

func (s *taskStorageImpl) GetAssignmentLog(ctx context.Context, c *domain.AssignmentLogCriteria) (*domain.AssignmentLogResponse, error) {

	response := &domain.AssignmentLogResponse{
		PagingResponse: &common.PagingResponse{
			Total: 0,
			Index: 0,
		},
		Logs: []*domain.AssignmentLog{},
	}

	selectClause := `*`

	query := s.c.Db.Instance.
		Table(`assignment_logs l`).
		Order(`l.start_time desc`)

	if c.StartTimeAfter != nil {
		query = query.Where(`l.start_time >= ?`, c.StartTimeAfter)
	}

	if c.StartTimeBefore != nil {
		query = query.Where(`l.start_time <= ?`, c.StartTimeBefore)
	}

	// paging
	var totalCount int64
	var offset int

	query.Count(&totalCount)

	if totalCount > int64(c.Size) {
		offset = (c.Index - 1) * c.Size
	}

	response.PagingResponse.Total = int(math.Ceil(float64(totalCount) / float64(c.Size)))
	response.PagingResponse.Index = c.Index

	query = query.Select(selectClause).Offset(offset).Limit(c.Size)

	rows, err := query.Rows()
	var logs []*assignmentLog
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		l := &assignmentLog{}
		_ = s.c.Db.Instance.ScanRows(rows, l)
		logs = append(logs, l)
	}
	response.Logs = s.toAssgnLogsDomain(logs)

	return response, nil
}
