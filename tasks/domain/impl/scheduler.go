package impl

import (
	"context"
	"github.com/go-co-op/gocron"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"gitlab.medzdrav.ru/prototype/tasks/logger"
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

func (r *reminderImpl) dueDateFunc(ctx context.Context, taskId string) {

	logger.L().Cmp("task-sch").Mth("due-date-fired").C(ctx).Dbg()

	r.handlerMutex.RLock()
	defer r.handlerMutex.RUnlock()

	if r.dueDateHandler != nil {
		r.dueDateHandler(ctx, taskId)
	}
}

func (r *reminderImpl) remindFunc(ctx context.Context, taskId string) {

	logger.L().Cmp("task-sch").Mth("remind-fired").C(ctx).Dbg()

	r.handlerMutex.RLock()
	defer r.handlerMutex.RUnlock()

	if r.reminderHandler != nil {
		r.reminderHandler(ctx, taskId)
	}
}

func (r *reminderImpl) ScheduleTask(ctx context.Context, ts *domain.Task) {

	l := logger.L().Cmp("task-sch").Mth("schedule").C(ctx).F(log.FF{"task": ts.Num})

	if r.config.IsFinalStatus(ctx, ts.Type, ts.Status) {
		return
	}

	if ts.DueDate != nil && ts.DueDate.After(time.Now().UTC()) {
		j, _ := r.dueDateScheduler.Every(1).Day().StartAt(*ts.DueDate).Do(r.dueDateFunc, ctx, ts.Id)
		l.TrcF("duedate scheduled at %v (next run: %v)\n", *ts.DueDate, j.NextRun())
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
				l.Err("not supported unit type %v", rmnd.BeforeDueDate.Unit)
				continue
			}

			d = d * int64(rmnd.BeforeDueDate.Value)
			schTime := ts.DueDate.Add(-time.Duration(d))

			if schTime.After(time.Now().UTC()) {
				_, _ = r.reminderScheduler.Every(1).Day().StartAt(schTime).Do(r.remindFunc, ctx, ts.Id)
				l.TrcF("scheduled reminder at %v", schTime)
			}

		}

		if rmnd.SpecificTime != nil && rmnd.SpecificTime.At != nil && rmnd.SpecificTime.At.After(time.Now().UTC()) {
			_, _ = r.reminderScheduler.Every(1).Day().StartAt(*rmnd.SpecificTime.At).Do(r.remindFunc, ctx, ts.Id)
			l.DbgF("scheduled reminder at %v", rmnd.SpecificTime.At)
		}

	}


}

func (r *reminderImpl) start(ctx context.Context) {

	l := logger.L().Cmp("task-sch").Mth("start").C(ctx)

	r.reminderScheduler.StartAsync()
	r.dueDateScheduler.StartAsync()

	l.Dbg("scheduler started")

	// TODO: retrieve tasks which have Reminder in the future
	// TODO: paging
	rs, err := r.storage.Search(ctx, &domain.SearchCriteria{
		PagingRequest: &common.PagingRequest{
			Size: 10000,
		},
		Status: &domain.Status{
			Status: domain.TS_OPEN,
		},
	})
	if err != nil {
		l.E(err).Err("task search failed")
		return
	}

	l.DbgF("found %d tasks\n", len(rs.Tasks))

	for _, t := range rs.Tasks {
		r.ScheduleTask(ctx, t)
	}
}

func (r *reminderImpl) StartAsync(ctx context.Context) {
	go r.start(ctx)
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
