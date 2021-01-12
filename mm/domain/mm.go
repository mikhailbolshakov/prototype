package domain

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/users"
	tasks2 "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/queue_model"
	"log"
)

type MMService interface {}

type serviceImpl struct {
	mmService mattermost.Service
	usersService users.Service
	tasksService tasks.Service
}

func NewMMService(mmService mattermost.Service, usersService users.Service, tasksService tasks.Service) MMService {

	s := &serviceImpl{
		mmService: mmService,
		usersService: usersService,
		tasksService: tasksService,
	}

	// setup handlers
	s.mmService.SetNewPostMessageHandler(s.postHandler)
	s.tasksService.SetTaskAssignedHandler(s.taskAssignedHandler)

	return s
}

func (s *serviceImpl) postHandler(post *mattermost.Post) {

	// get user by MM user id
	user := s.usersService.GetByMMId(post.UserId)

	if user != nil && user.MMChannelId == post.ChannelId {

		// retrieves tasks by channel
		ts := s.tasksService.GetByChannelId(user.MMChannelId)

		// check if there is open task
		newTask := true
		if len(ts) > 0 {

			for _, t := range ts {
				// TODO: it's simplification
				// a correct check should verify if there are no tasks with close time > post time
				// otherwise this post relates to the closed task and somehow hasn't been delivered in time
				if t.Status.Status == "open" {
					newTask = false
					break
				}
			}

		}

		if newTask {

			postTime := kit.TimeFromMillis(post.CreateAt)
			ts, _ := ptypes.TimestampProto(postTime)

			// create a new task
			createdTask, err := s.tasksService.CreateTask(&tasks2.NewTaskRequest{
				Type:        &tasks2.Type{
					Type:    "client",
					Subtype: "medical-request",
				},
				ReportedBy:  user.Username,
				ReportedAt:  ts,
				Description: "Обращение клиента",
				Title:       "Обращение клиента",
				DueDate:     nil,
				Assignee:    &tasks2.Assignee{
					Group: "consultant",
				},
				ChannelId:   user.MMChannelId,
			})
			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("task created: %s", createdTask.Id)

			if err := s.tasksService.MakeTransition(&tasks2.MakeTransitionRequest{
				TaskId:       createdTask.Id,
				TransitionId: "2",
			}); err != nil {
				log.Println(err)
				return
			}

			if err := s.sendTriggerPost(TP_CLIENT_NEW_REQUEST, user.MMId, user.MMChannelId, triggerPostParams{
				"client.name": fmt.Sprintf("%s", user.FirstName),
			}); err != nil {
				log.Println(err)
				return
			}

		} else {
			log.Println("task found")
		}

	}

}

func (s *serviceImpl) taskAssignedHandler(task *queue_model.Task) {

	if task.Type.Type == "client" && task.Type.SubType == "medical-request" && task.Assignee.User != "" {

		assigneeUser := s.usersService.GetByUsername(task.Assignee.User)

		if _, err := s.mmService.SubscribeUser(&mattermost.SubscribeUserRequest{
			UserId:    assigneeUser.MMId,
			ChannelId: task.ChannelId,
		}); err != nil {
			log.Println(err)
			return
		}

		reportedUser := s.usersService.GetByUsername(task.Reported.By)

		if err := s.sendTriggerPost(TP_CLIENT_REQUEST_ASSIGNED, reportedUser.MMId, task.ChannelId, triggerPostParams{
			"consultant.first-name": assigneeUser.FirstName,
			"consultant.last-name": assigneeUser.LastName,
			"consultant.url": "https://prodoctorov.ru/media/photo/tula/doctorimage/589564/432638-589564-ezhikov_l.jpg",
		}); err != nil {
			log.Println(err)
			return
		}

		if err := s.sendTriggerPost(TP_CONSULTANT_REQUEST_ASSIGNED, assigneeUser.MMId, task.ChannelId, triggerPostParams{
			"client.first-name": reportedUser.FirstName,
			"client.last-name": reportedUser.LastName,
			"client.phone": reportedUser.Phone,
			"client.url": "https://cdn5.vedomosti.ru/crop/image/2020/2s/qmb9n/original-yi0.jpg?height=934&width=1660",
			"client.med-card": "https://pmed.moi-service.ru/profile/medcard",
		}); err != nil {
			log.Println(err)
			return
		}


	}

}