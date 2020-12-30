package storage

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/tasks/infrastructure"
	"math"
	"time"
)

type TaskStorage interface {
	Create(t *Task) (*Task, error)
	Get(id string) *Task
	Update(t *Task) (*Task, error)
	GetByChannel(channelId string) []*Task
	CreateHistory(h *History) (*History, error)
	Search(cr *SearchCriteria) (*SearchResponse, error)
}

type taskStorageImpl struct {
	infr *infrastructure.Container
}

func NewStorage(infr *infrastructure.Container) TaskStorage {
	return &taskStorageImpl{infr: infr}
}

func (s *taskStorageImpl) Create(task *Task) (*Task, error) {

	t := time.Now()
	task.CreatedAt, task.UpdatedAt = t, t

	result := s.infr.Db.Instance.Create(task)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO: put to Redis
	// TODO: save history

	return task, nil
}

func (s *taskStorageImpl) Get(id string) *Task {

	// TODO: get from Redis

	task := &Task{Id: id}
	s.infr.Db.Instance.First(task)
	return task
}

func (s *taskStorageImpl) Update(task *Task) (*Task, error) {

	task.UpdatedAt = time.Now()

	result := s.infr.Db.Instance.Save(task)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO: put to Redis
	// TODO: save history

	return task, nil
}

func (s *taskStorageImpl) GetByChannel(channelId string) []*Task {
	var tasks []*Task
	s.infr.Db.Instance.Where("channel_id = ?", channelId).Find(&tasks)
	return tasks
}

func (s *taskStorageImpl) CreateHistory(h *History) (*History, error) {
	result := s.infr.Db.Instance.Create(h)
	if result.Error != nil {
		return nil, result.Error
	}
	return h, nil
}

func (s *taskStorageImpl) Search(cr *SearchCriteria) (*SearchResponse, error) {

	response := &SearchResponse{
		PagingResponse: &common.PagingResponse{
			Total: 0,
			Index: 0,
		},
		Tasks: []*Task{},
	}

	selectClause := `*`

	query := s.infr.Db.Instance.
		Table(`tasks t`).
		Where(`t.deleted_at is null`)

	if cr.Num != "" {
		query = query.Where(`t.num = ?`, cr.Num)
	}

	if cr.Type != "" {
		query = query.Where(`t.type = ?`, cr.Type)
	}

	if cr.SubType != "" {
		query = query.Where(`t.subtype = ?`, cr.SubType)
	}

	if cr.Status != "" {
		query = query.Where(`t.status = ?`, cr.Status)
	}

	if cr.SubStatus != "" {
		query = query.Where(`t.substatus = ?`, cr.SubStatus)
	}

	if cr.AssigneeUser != "" {
		query = query.Where(`t.assignee_user = ?`, cr.AssigneeUser)
	}

	if cr.AssigneeGroup != "" {
		query = query.Where(`t.assignee_group = ?`, cr.AssigneeGroup)
	}

	// paging
	var totalCount int64
	var offset int

	query.Count(&totalCount)

	if totalCount > int64(cr.Size) {
		offset = (cr.Index - 1) * cr.Size
	}

	response.PagingResponse.Total = int(math.Ceil(float64(totalCount) / float64(cr.Size)))
	response.PagingResponse.Index = cr.Index

	query = query.Select(selectClause).Offset(offset).Limit(cr.Size)

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &Task{}
		_ = s.infr.Db.Instance.ScanRows(rows, task)
		response.Tasks = append(response.Tasks, task)
	}

	return response, nil
}

