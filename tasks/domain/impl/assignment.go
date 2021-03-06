package impl

import (
	"context"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/log"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"gitlab.medzdrav.ru/prototype/tasks/logger"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type assignmentTask struct {
	taskType *domain.Type
	cfg      *domain.Config
	quit     chan struct{}
	cancel   context.CancelFunc
	ctx      context.Context
	run      *atomic.Bool
	sync.Mutex
}

func (t *assignmentTask) setRun(v bool) {
	t.run.Store(v)
}

func (t *assignmentTask) getRun() bool {
	return t.run.Load()
}

func NewAssignmentDaemon(cfgService domain.ConfigService,
	taskService domain.TaskService,
	userSearchService domain.UserService,
	storage domain.TaskStorage) domain.AssignmentDaemon {
	d := &daemonImpl{
		taskTypes:   []*assignmentTask{},
		cfgService:  cfgService,
		taskService: taskService,
		userService: userSearchService,
		storage:     storage,
	}
	return d
}

type daemonImpl struct {
	taskTypes   []*assignmentTask
	cfgService  domain.ConfigService
	taskService domain.TaskService
	userService domain.UserService
	storage     domain.TaskStorage
}

func (d *daemonImpl) l() log.CLogger {
	return logger.L().Cmp("assign-daemon")
}

func (d *daemonImpl) assign(ctx context.Context, tt *assignmentTask) error {

	ll := d.l().Mth("assign").C(ctx)

	logSuccess := func(log *domain.AssignmentLog) {
		log.Status = "success"
		t := time.Now().UTC()
		log.FinishTime = &t
		_, _ = d.storage.SaveAssignmentLog(ctx, log)
	}

	logFail := func(log *domain.AssignmentLog, err error) {
		log.Status = "fail"
		log.Error = err.Error()
		t := time.Now().UTC()
		log.FinishTime = &t
		_, _ = d.storage.SaveAssignmentLog(ctx, log)
	}

	for _, rule := range tt.cfg.AssignmentRules {

		l, _ := d.storage.SaveAssignmentLog(ctx, &domain.AssignmentLog{
			StartTime:       time.Now().UTC(),
			Status:          "in-progress",
			RuleCode:        rule.Code,
			RuleDescription: rule.Description,
			UsersInPool:     0,
			TasksToAssign:   0,
			Assigned:        0,
		})

		// search users for the assignee pool
		usersRs, err := d.userService.Search(ctx, &pb.SearchRequest{
			Paging: &pb.PagingRequest{
				Size:  100,
				Index: 1,
			},
			UserType:       rule.UserPool.Type,
			UserGroup:      rule.UserPool.Group,
			Status:         "active",
			OnlineStatuses: rule.UserPool.Statuses,
		})
		if err != nil {
			logFail(l, err)
			ll.E(err).Err("user search failed")
			return err
		}

		l.UsersInPool = len(usersRs.Users)

		// search task to be assigned by the rule
		cr := &domain.SearchCriteria{
			PagingRequest: &common.PagingRequest{
				Size:  100,
				Index: 1,
			},
			Status:   rule.Source.Status,
			Assignee: rule.Source.Assignee,
			Type:     tt.taskType,
		}
		rs, err := d.taskService.Search(ctx, cr)
		if err != nil {
			logFail(l, err)
			ll.E(err).Err("task search failed")
			return err
		}

		l.TasksToAssign = len(rs.Tasks)

		ll.F(log.FF{"rule": rule.Code}).Dbg("rule fired").TrcF("users=%d tasks=%d", l.UsersInPool, l.TasksToAssign)

		if len(rs.Tasks) == 0 {
			logSuccess(l)
			break
		}

		usersPool := map[string]*pb.User{}
		for _, u := range usersRs.Users {
			usersPool[u.Username] = u
		}

		userToAssign := &pb.User{}

		// go through tasks
		for _, t := range rs.Tasks {

			if len(usersPool) == 0 {
				// TODO: message from here to notify task there are no available users to assign
				ll.F(log.FF{"task": t.Num}).Dbg("no users")
				break
			}

			// take the first user
			for _, u := range usersPool {
				userToAssign = u
				break
			}

			// delete it from the pool
			delete(usersPool, userToAssign.Username)

			// assign task to user
			t, err = d.taskService.SetAssignee(ctx, t.Id, &domain.Assignee{
				UserId: userToAssign.Id,
			})
			if err != nil {
				logFail(l, err)
				ll.E(err).Err("set assignee failed")
				return err
			}
			ll.F(log.FF{"task": t.Num}).Inf("assigned user=%s", userToAssign.Username)

			// if rule specifies a target status then make transition
			if rule.Target != nil && rule.Target.Status != nil {

				tr, err := d.cfgService.FindTransition(ctx, t.Type, t.Status, rule.Target.Status)
				if err != nil {
					logFail(l, err)
					ll.E(err).Err("find transition failed")
					return err
				}
				t, err = d.taskService.MakeTransition(ctx, t.Id, tr.Id)
				if err != nil {
					logFail(l, err)
					ll.E(err).Err("make transition failed")
					return err
				}

			}
			l.Assigned++

		}

		logSuccess(l)
	}

	return nil

}

func (d *daemonImpl) Run(ctx context.Context) {

	for _, t := range d.taskTypes {

		t.setRun(true)

		go func(tt *assignmentTask) {

			defer tt.setRun(false)

			for {

				select {
				// TODO: configuration
				case <-time.Tick(time.Second * 20):
					if err := d.assign(ctx, tt); err != nil {
						return
					}
				case <-tt.ctx.Done():
					d.l().Mth("run").F(log.FF{"task-type": tt.taskType}).C(ctx).Dbg("cancelled")
					return
				}

			}

		}(t)

	}

}

func (d *daemonImpl) Stop(ctx context.Context) error {

	for _, t := range d.taskTypes {
		if t.getRun() {
			t.cancel()
		}
	}

	return nil
}

func (d *daemonImpl) Init(ctx context.Context) error {

	d.taskTypes = []*assignmentTask{}

	for _, cfg := range d.cfgService.GetAll(ctx) {

		if cfg.AssignmentRules != nil && len(cfg.AssignmentRules) > 0 {

			ctx, cancel := context.WithCancel(ctx)

			d.taskTypes = append(d.taskTypes, &assignmentTask{
				taskType: cfg.Type,
				cfg:      cfg,
				quit:     make(chan struct{}),
				ctx:      ctx,
				cancel:   cancel,
				run:      atomic.NewBool(false),
			})
		}

	}

	return nil

}
