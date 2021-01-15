package domain

import (
	"github.com/go-co-op/gocron"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
	"log"
	"sync"
	"time"
)

type TaskSchedulerHandler func(taskId string)

type TaskScheduler interface {
	SetReminderHandler(h TaskSchedulerHandler)
	SetDueDateHandler(h TaskSchedulerHandler)
	StartAsync()
	ScheduleTask(ts *Task)
}

type reminderImpl struct {
	storage           storage.TaskStorage
	config            ConfigService
	reminderScheduler *gocron.Scheduler
	dueDateScheduler  *gocron.Scheduler
	reminderHandler   TaskSchedulerHandler
	dueDateHandler    TaskSchedulerHandler
	handlerMutex      sync.RWMutex
}

func NewScheduler(config ConfigService, storage storage.TaskStorage) TaskScheduler {
	return &reminderImpl{
		storage:           storage,
		config:            config,
		reminderScheduler: gocron.NewScheduler(time.UTC),
		dueDateScheduler:  gocron.NewScheduler(time.UTC),
		handlerMutex:      sync.RWMutex{},
	}
}

func (r *reminderImpl) dueDateFunc(taskId string) {

	r.handlerMutex.RLock()
	defer r.handlerMutex.RUnlock()

	if r.dueDateHandler != nil {
		r.dueDateHandler(taskId)
	}
}

func (r *reminderImpl) remindFunc(taskId string) {

	r.handlerMutex.RLock()
	defer r.handlerMutex.RUnlock()

	if r.reminderHandler != nil {
		r.reminderHandler(taskId)
	}
}

func (r *reminderImpl) ScheduleTask(ts *Task) {

	if r.config.IsFinalStatus(ts.Type, ts.Status) {
		return
	}

	if ts.DueDate != nil && ts.DueDate.After(time.Now().UTC()) {
		_, _ = r.dueDateScheduler.Every(1).Hours().StartAt(*ts.DueDate).Do(r.dueDateFunc, ts.Id)
		log.Printf("scheduler set for due date: %v", ts.DueDate)
	}

	for _, rmnd := range ts.Reminders {

		if rmnd.BeforeDueDate != nil && ts.DueDate != nil && rmnd.BeforeDueDate.Unit != "" {

			var d int64

			switch rmnd.BeforeDueDate.Unit {
			case seconds:
				d = int64(time.Second)
			case minutes:
				d = int64(time.Minute)
			case hours:
				d = int64(time.Hour)
			case days:
				d = 24 * int64(time.Hour)
			default:
				log.Println("ERROR: not supported unit type ", rmnd.BeforeDueDate.Unit)
				continue
			}

			d = d * int64(rmnd.BeforeDueDate.Value)
			schTime := ts.DueDate.Add(-time.Duration(d))

			if schTime.After(time.Now().UTC()) {
				_, _ = r.reminderScheduler.Every(1).Hours().StartAt(schTime).Do(r.remindFunc, ts.Id)
				log.Printf("scheduler set for remind 'before due date': %v", schTime)
			}

		}

		if rmnd.SpecificTime != nil && rmnd.SpecificTime.At != nil {
			_, _ = r.reminderScheduler.Every(1).Day().StartAt(*rmnd.SpecificTime.At).Do(r.remindFunc, ts.Id)
			log.Printf("scheduler set for remind 'specific time': %v", rmnd.SpecificTime.At)
		}

	}


}

func (r *reminderImpl) start() {

	rs, err := r.storage.Search(&storage.SearchCriteria{
		PagingRequest: &common.PagingRequest{
			Size: 10000,
		},
	})
	if err != nil {
		log.Printf("ERROR: start reminder, search error: %v", err)
		return
	}

	log.Printf("preparing reminders... found %d tasks", len(rs.Tasks))

	for _, dto := range rs.Tasks {
		r.ScheduleTask(fromDto(dto))
	}

	r.reminderScheduler.StartAsync()
	r.dueDateScheduler.StartAsync()
}

func (r *reminderImpl) StartAsync() {
	go r.start()
}

func (r *reminderImpl) SetReminderHandler(h TaskSchedulerHandler) {
	r.handlerMutex.Lock()
	defer r.handlerMutex.Unlock()
	r.reminderHandler = h
}

func (r *reminderImpl) SetDueDateHandler(h TaskSchedulerHandler) {
	r.handlerMutex.Lock()
	defer r.handlerMutex.Unlock()
	r.dueDateHandler = h
}
