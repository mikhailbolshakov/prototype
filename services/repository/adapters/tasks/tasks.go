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
	MakeTransition(rq *pb.MakeTransitionRequest) error
	SetTaskDueDateHandler(h TaskHandler)
}

type TaskHandler func(task *queue_model.Task)

type serviceImpl struct {
	queue               queue.Queue
	taskDueDateHandler	TaskHandler
	pb.TasksClient
}

func newImpl(queue queue.Queue) *serviceImpl {
	a := &serviceImpl{
		queue: queue,
	}
	return a
}

func (u *serviceImpl) SetTaskDueDateHandler(h TaskHandler) {
	u.taskDueDateHandler = h
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

	taskDueDateChan := make(chan []byte)
	err := u.queue.Subscribe("tasks.duedate", taskDueDateChan)
	if err != nil {
		return err
	}

	go func() {

		for {
			select {
			case msg := <-taskDueDateChan:

				task := &queue_model.Task{}
				_ = json.Unmarshal(msg, task)

				log.Printf("assigned task event received %v", task)

				if u.taskDueDateHandler != nil {
					u.taskDueDateHandler(task)
				}

			}
		}

	}()

	return nil
}
