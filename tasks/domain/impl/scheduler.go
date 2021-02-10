package impl

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"sync"
	"time"
)

type reminderImpl struct {
	storage           domain.TaskStorage
	config            domain.ConfigService
	reminderScheduler *gocron.Scheduler
	dueDateScheduler  *gocron.Scheduler
	reminderHandler   domain.TaskSchedulerHandler
	dueDateHandler    domain.TaskSchedulerHandler
	handlerMutex      sync.RWMutex
}

func NewScheduler(config domain.ConfigService, storage domain.TaskStorage) domain.TaskScheduler {
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

func (r *reminderImpl) ScheduleTask(ts *domain.Task) {

	if r.config.IsFinalStatus(ts.Type, ts.Status) {
		return
	}

	if ts.DueDate != nil && ts.DueDate.After(time.Now().UTC()) {
		_, _ = r.dueDateScheduler.Every(1).Hours().StartAt(*ts.DueDate).Do(r.dueDateFunc, ts.Id)
		log.DbgF("scheduler set for due date: %v", ts.DueDate)
	}

	for _, rmnd := range ts.Reminders {

		if rmnd.BeforeDueDate != nil && ts.DueDate != nil && rmnd.BeforeDueDate.Unit != "" {

			var d int64

			switch rmnd.BeforeDueDate.Unit {
			case "seconds":
				d = int64(time.Second)
			case "minutes":
				d = int64(time.Minute)
			case "hours":
				d = int64(time.Hour)
			case "days":
				d = 24 * int64(time.Hour)
			default:
				log.Err(fmt.Errorf("ERROR: not supported unit type %v", rmnd.BeforeDueDate.Unit), true)
				continue
			}

			d = d * int64(rmnd.BeforeDueDate.Value)
			schTime := ts.DueDate.Add(-time.Duration(d))

			if schTime.After(time.Now().UTC()) {
				_, _ = r.reminderScheduler.Every(1).Hours().StartAt(schTime).Do(r.remindFunc, ts.Id)
				log.DbgF("scheduler set for remind 'before due date': %v", schTime)
			}

		}

		if rmnd.SpecificTime != nil && rmnd.SpecificTime.At != nil && rmnd.SpecificTime.At.After(time.Now().UTC()) {
			_, _ = r.reminderScheduler.Every(1).Day().StartAt(*rmnd.SpecificTime.At).Do(r.remindFunc, ts.Id)
			log.DbgF("scheduler set for remind 'specific time': %v", rmnd.SpecificTime.At)
		}

	}


}

func (r *reminderImpl) start() {

	// TODO: retrieve tasks which have Reminder in the future
	rs, err := r.storage.Search(&domain.SearchCriteria{
		PagingRequest: &common.PagingRequest{
			Size: 10000,
		},
		Status: &domain.Status{
			Status: domain.TS_OPEN,
		},
	})
	if err != nil {
		log.Err(err, true)
		return
	}

	log.DbgF("preparing reminders... found %d tasks", len(rs.Tasks))

	for _, t := range rs.Tasks {
		r.ScheduleTask(t)
	}

	r.reminderScheduler.StartAsync()
	r.dueDateScheduler.StartAsync()
}

func (r *reminderImpl) StartAsync() {
	go r.start()
}

func (r *reminderImpl) SetReminderHandler(h domain.TaskSchedulerHandler) {
	r.handlerMutex.Lock()
	defer r.handlerMutex.Unlock()
	r.reminderHandler = h
}

func (r *reminderImpl) SetDueDateHandler(h domain.TaskSchedulerHandler) {
	r.handlerMutex.Lock()
	defer r.handlerMutex.Unlock()
	r.dueDateHandler = h
}
