package repository

import "time"

type TaskStorage interface {
	Create(t *Task) (*Task, error)
	Get(id string) *Task
	Update(t *Task) (*Task, error)
}

type TaskStorageImpl struct {}

func NewStorage() TaskStorage {
	return &TaskStorageImpl{}
}

func (s *TaskStorageImpl) Create(task *Task) (*Task, error) {

	t := time.Now()
	task.CreatedAt, task.UpdatedAt = t, t

	result := storage.Instance.Create(task)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO: put to Redis
	// TODO: save history

	return task, nil
}

func (s *TaskStorageImpl) Get(id string) *Task {

	// TODO: get from Redis

	task := &Task{Id: id}
	storage.Instance.First(task)
	return task
}

func (s *TaskStorageImpl) Update(task *Task) (*Task, error) {

	task.UpdatedAt = time.Now()

	result := storage.Instance.Save(task)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO: put to Redis
	// TODO: save history

	return task, nil
}