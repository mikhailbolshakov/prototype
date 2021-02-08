package domain

import (
	"context"
	"gitlab.medzdrav.ru/prototype/kit/common"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"gitlab.medzdrav.ru/prototype/tasks/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
	"log"
	"sync"
	"time"
)

type AssignmentDaemon interface {
	Run()
	Stop() error
	Init() error
}

type assignmentTask struct {
	taskType *Type
	cfg      *Config
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


func NewAssignmentDaemon(cfgService ConfigService,
	taskService TaskService,
	searchService TaskSearchService,
	userSearchService users.Service,
	storage storage.TaskStorage) AssignmentDaemon {
	d := &daemonImpl{
		taskTypes:     []*assignmentTask{},
		cfgService:    cfgService,
		taskService:   taskService,
		searchService: searchService,
		userService:   userSearchService,
		storage:       storage,
	}
	return d
}

type daemonImpl struct {
	taskTypes     []*assignmentTask
	cfgService    ConfigService
	taskService   TaskService
	searchService TaskSearchService
	userService   users.Service
	storage       storage.TaskStorage
}

func (d *daemonImpl) assign(tt *assignmentTask) error {

	logSuccess := func(log *storage.AssignmentLog) {
		log.Status = "success"
		t := time.Now().UTC()
		log.FinishTime = &t
		_, _ = d.storage.SaveAssignmentLog(log)
	}

	logFail := func(log *storage.AssignmentLog, err error) {
		log.Status = "fail"
		log.Error = err.Error()
		t := time.Now().UTC()
		log.FinishTime = &t
		_, _ = d.storage.SaveAssignmentLog(log)
	}

	for _, rule := range tt.cfg.AssignmentRules {

		l, _ := d.storage.SaveAssignmentLog(&storage.AssignmentLog{
			StartTime:       time.Now().UTC(),
			Status:          "in-progress",
			RuleCode:        rule.Code,
			RuleDescription: rule.Description,
			UsersInPool:     0,
			TasksToAssign:   0,
			Assigned:        0,
		})

		log.Printf("assignment rule is fired %v", rule)

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
			logSuccess(l)
			log.Fatal(err)
			return err
		}

		l.UsersInPool = len(usersRs.Users)

		// search task to be assigned by the rule
		cr := &SearchCriteria{
			PagingRequest: &common.PagingRequest{
				Size:  100,
				Index: 1,
			},
			Status:   rule.Source.Status,
			Assignee: rule.Source.Assignee,
			Type:     tt.taskType,
		}
		rs, err := d.searchService.Search(cr)
		if err != nil {
			logFail(l, err)
			log.Fatal(err)
			return err
		}

		l.TasksToAssign = len(rs.Tasks)

		if len(rs.Tasks) == 0 {
			logSuccess(l)
			log.Printf("no task to assign found")
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
				log.Printf("no users available for task %s", t.Id)
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
			t, err = d.taskService.SetAssignee(t.Id, &Assignee{
				UserId: userToAssign.Id,
			})
			if err != nil {
				logFail(l, err)
				log.Fatal(err)
				return err
			}
			log.Printf("user %s assigned on task %s", userToAssign.Username, t.Id)

			// if rule specifies a target status then make transition
			if rule.Target != nil && rule.Target.Status != nil {

				tr, err := d.cfgService.FindTransition(t.Type, t.Status, rule.Target.Status)
				if err != nil {
					logFail(l, err)
					log.Fatal(err)
					return err
				}
				t, err = d.taskService.MakeTransition(t.Id, tr.Id)
				if err != nil {
					logFail(l, err)
					log.Fatal(err)
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
				//case <-tt.quit:
				//	log.Printf("assignment task of type %v is cancelled", tt.taskType)
				//	return
				case <-tt.ctx.Done():
					log.Printf("assignment task of type %v is cancelled", tt.taskType)
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
			//t.quit <- struct{}{}
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
