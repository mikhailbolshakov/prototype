package impl

import (
	"context"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/log"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"sync"
	"time"
)

type assignmentTask struct {
	taskType *domain.Type
	cfg      *domain.Config
	quit     chan struct{}
	cancel   context.CancelFunc
	ctx      context.Context
	run      bool
	sync.Mutex
}

func (t *assignmentTask) setRun(v bool) {
	t.Lock()
	defer t.Unlock()
	t.run = v
}

func (t *assignmentTask) getRun() bool {
	t.Lock()
	defer t.Unlock()
	return t.run
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

func (d *daemonImpl) assign(tt *assignmentTask) error {

	logSuccess := func(log *domain.AssignmentLog) {
		log.Status = "success"
		t := time.Now().UTC()
		log.FinishTime = &t
		_, _ = d.storage.SaveAssignmentLog(log)
	}

	logFail := func(log *domain.AssignmentLog, err error) {
		log.Status = "fail"
		log.Error = err.Error()
		t := time.Now().UTC()
		log.FinishTime = &t
		_, _ = d.storage.SaveAssignmentLog(log)
	}

	for _, rule := range tt.cfg.AssignmentRules {

		l, _ := d.storage.SaveAssignmentLog(&domain.AssignmentLog{
			StartTime:       time.Now().UTC(),
			Status:          "in-progress",
			RuleCode:        rule.Code,
			RuleDescription: rule.Description,
			UsersInPool:     0,
			TasksToAssign:   0,
			Assigned:        0,
		})

		log.DbgF("assignment rule is fired %v", rule)

		// search users for the assignee pool
		usersRs, err := d.userService.Search(&pb.SearchRequest{
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
			log.Err(err, true)
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
		rs, err := d.taskService.Search(cr)
		if err != nil {
			logFail(l, err)
			log.Err(err, true)
			return err
		}

		l.TasksToAssign = len(rs.Tasks)

		if len(rs.Tasks) == 0 {
			logSuccess(l)
			log.Dbg("no task to assign found")
			break
		}

		usersPool := map[string]*pb.User{}
		for _, u := range usersRs.Users {
			usersPool[u.Username] = u
		}

		var userToAssign *pb.User

		// go through tasks
		for _, t := range rs.Tasks {

			if len(usersPool) == 0 {
				// TODO: message from here to notify task there are no available users to assign
				log.DbgF("no users available for task %s", t.Id)
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
			t, err = d.taskService.SetAssignee(t.Id, &domain.Assignee{
				UserId: userToAssign.Id,
			})
			if err != nil {
				logFail(l, err)
				log.Err(err, true)
				return err
			}
			log.DbgF("user %s assigned on task %s", userToAssign.Username, t.Id)

			// if rule specifies a target status then make transition
			if rule.Target != nil && rule.Target.Status != nil {

				tr, err := d.cfgService.FindTransition(t.Type, t.Status, rule.Target.Status)
				if err != nil {
					logFail(l, err)
					log.Err(err, true)
					return err
				}
				t, err = d.taskService.MakeTransition(t.Id, tr.Id)
				if err != nil {
					logFail(l, err)
					log.Err(err, true)
					return err
				}

			}
			l.Assigned++

		}

		logSuccess(l)
	}

	return nil

}

func (d *daemonImpl) Run() {

	for _, t := range d.taskTypes {

		t.setRun(true)

		go func(tt *assignmentTask) {

			defer tt.setRun(false)

			for {

				select {
				// TODO: configuration
				case <-time.Tick(time.Second * 20):
					if err := d.assign(tt); err != nil {
						return
					}
				case <-tt.ctx.Done():
					log.DbgF("assignment task of type %v is cancelled", tt.taskType)
					return
				}

			}

		}(t)

	}

}

func (d *daemonImpl) Stop() error {

	for _, t := range d.taskTypes {
		if t.getRun() {
			t.cancel()
		}
	}

	return nil
}

func (d *daemonImpl) Init() error {

	d.taskTypes = []*assignmentTask{}

	for _, cfg := range d.cfgService.GetAll() {

		if cfg.AssignmentRules != nil && len(cfg.AssignmentRules) > 0 {

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)

			d.taskTypes = append(d.taskTypes, &assignmentTask{
				taskType: cfg.Type,
				cfg:      cfg,
				quit:     make(chan struct{}),
				ctx:      ctx,
				cancel:   cancel,
			})
		}

	}

	return nil

}
