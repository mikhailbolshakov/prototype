package tasks

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/queue_model"
	"log"
)

type Service interface {
	GetByChannelId(channelId string) []*pb.Task
	CreateTask(rq *pb.NewTaskRequest) (*pb.Task, error)
	SetTaskAssignedHandler(h TaskHandler)
	SetTaskRemindHandler(h TaskHandler)
	MakeTransition(rq *pb.MakeTransitionRequest) error
}

type serviceImpl struct {
	queue               queue.Queue
	taskAssignedHandler TaskHandler
	taskRemindHandler   TaskHandler
	pb.TasksClient
}

type TaskHandler func(task *queue_model.Task)

func newImpl(queue queue.Queue) *serviceImpl {
	a := &serviceImpl{
		queue: queue,
	}
	return a
}

func (u *serviceImpl) SetTaskAssignedHandler(h TaskHandler) {
	u.taskAssignedHandler = h
}

func (u *serviceImpl) SetTaskRemindHandler(h TaskHandler) {
	u.taskRemindHandler = h
}

func (u *serviceImpl) GetByChannelId(channelId string) []*pb.Task {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs, err := u.GetByChannel(ctx, &pb.GetByChannelRequest{ChannelId: channelId})
	if err != nil {
		log.Printf("error: %v", err)
		return []*pb.Task{}
	}
	return rs.Tasks
}

func (u *serviceImpl) MakeTransition(rq *pb.MakeTransitionRequest) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := u.TasksClient.MakeTransition(ctx, rq)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	return nil
}

func (u *serviceImpl) CreateTask(rq *pb.NewTaskRequest) (*pb.Task, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return u.New(ctx, rq)
}

func (u *serviceImpl) listenTaskQueue() error {

	taskAssignedChan := make(chan []byte)
	err := u.queue.Subscribe("tasks.assigned", taskAssignedChan)
	if err != nil {
		return err
	}

	taskRemindChan := make(chan []byte)
	err = u.queue.Subscribe("tasks.remind", taskRemindChan)
	if err != nil {
		return err
	}

	go func() {

		for {
			select {
			case msg := <-taskAssignedChan:

				task := &queue_model.Task{}
				_ = json.Unmarshal(msg, task)

				log.Printf("assigned task event received %v", task)

				if u.taskAssignedHandler != nil {
					u.taskAssignedHandler(task)
				}

			case msg := <-taskRemindChan:

				task := &queue_model.Task{}
				_ = json.Unmarshal(msg, task)

				log.Printf("remind task event received %v", task)

				if u.taskRemindHandler != nil {
					u.taskRemindHandler(task)
				}
			}
		}

	}()

	return nil
}
