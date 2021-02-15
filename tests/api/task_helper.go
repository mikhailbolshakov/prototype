package api

import (
	"context"
	"encoding/json"
	"fmt"
	taskApi "gitlab.medzdrav.ru/prototype/api/public/tasks"
	"time"
)

func (h *TestHelper) CreateTask(rq []byte) (*taskApi.Task, error) {

	rs, err := h.POST(fmt.Sprintf("%s/api/tasks", BASE_URL), rq)
	if err != nil {
		return nil, err
	} else {

		var task *taskApi.Task
		err = json.Unmarshal(rs, &task)
		if err != nil {
			return nil, err
		}

		if task == nil || task.Id == "" {
			return nil, fmt.Errorf("task invalid")
		}

		return task, nil
	}

}

func (h *TestHelper) MakeTransition(taskId, transition string) (*taskApi.Task, error) {

	rs, err := h.POST(fmt.Sprintf("%s/api/tasks/%s/transitions/%s", BASE_URL, taskId, transition), []byte{})
	if err != nil {
		return nil, err
	} else {

		var task *taskApi.Task
		err = json.Unmarshal(rs, &task)
		if err != nil {
			return nil, err
		}

		if task == nil || task.Id == "" {
			return nil, fmt.Errorf("task invalid")
		}

		return task, nil
	}
}

func (h *TestHelper) GetTask(taskId string) (*taskApi.Task, error) {

	rs, err := h.GET(fmt.Sprintf("%s/api/tasks/%s", BASE_URL, taskId))
	if err != nil {
		return nil, err
	}

	var user *taskApi.Task
	err = json.Unmarshal(rs, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *TestHelper) SearchByChannel(channelId string) ([]*taskApi.Task, error) {

	rs, err := h.GET(fmt.Sprintf("%s/api/tasks?channel=%s", BASE_URL, channelId))
	if err != nil {
		return nil, err
	}

	var sr *taskApi.SearchResponse
	err = json.Unmarshal(rs, &sr)
	if err != nil {
		return nil, err
	}

	return sr.Tasks, nil
}

func (h *TestHelper) TaskAwaitByChannel(channelId string, timeout time.Duration) (*taskApi.Task, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {

		case <- time.After(time.Second * 5):
			tasks, err := h.SearchByChannel(channelId)
			if err != nil {
				return nil, err
			}
			if len(tasks) > 0 {
				return tasks[0], nil
			}
			fmt.Printf("search by channel %s, not found\n", channelId)
		case <- ctx.Done():
			return nil, fmt.Errorf("timeout")
		}
	}
}

func (h *TestHelper) AssignTaskAwait(taskId string, timeout time.Duration) (*taskApi.Task, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {

		case <- time.After(time.Second * 5):
			task, err := h.GetTask(taskId)
			if err != nil {
				return nil, err
			}
			fmt.Printf("task retrieved, status %s\n", task.Status.SubStatus)
			if task.Status.SubStatus == "assigned" {
				return task, nil
			}
		case <- ctx.Done():
			return nil, fmt.Errorf("timeout")
		}
	}
}