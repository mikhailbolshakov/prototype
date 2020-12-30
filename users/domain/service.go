package domain

import (
	"github.com/golang/protobuf/ptypes"
	"gitlab.medzdrav.ru/prototype/kit"
	tasks2 "gitlab.medzdrav.ru/prototype/proto/tasks"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
	"log"
)

type UserService interface {
	Create(user *User) (*User, error)
	GetByUsername(username string) *User
	GetByMMId(mmId string) *User
}

type userServiceImpl struct {
	storage      storage.UserStorage
	muttermost   mattermost.Service
	tasksService tasks.Service
}

func NewUserService(storage storage.UserStorage, mmService mattermost.Service, tasksService tasks.Service) UserService {
	s := &userServiceImpl{
		storage:    storage,
		muttermost: mmService,
		tasksService: tasksService,
	}

	s.muttermost.SetNewPostMessageHandler(s.postHandler)

	return s
}

// TODO: this shouldn't be here
func (u *userServiceImpl) postHandler(post *mattermost.MMPost) {

	// get user by MM user id
	user := u.storage.GetByMMId(post.UserId)

	if user != nil && user.MMChannelId == post.ChannelId {

		// retrieves tasks by channel
		ts := u.tasksService.GetByChannelId(user.MMChannelId)

		// check if there is open task
		newTask := true
		if len(ts) > 0 {

			for _, t := range ts {
				// TODO: it's simplification
				// the correct check should verifies there are no tasks with close time > post time
				// otherwise this post relates to the closed task
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
			createdTask, err := u.tasksService.CreateTask(&tasks2.NewTaskRequest{
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
				log.Print(err)
				return
			}

			log.Printf("task created: %s", createdTask.Id)

		} else {
			log.Println("task found")
		}

	}

}

func (u *userServiceImpl) Create(user *User) (*User, error) {

	// create a new user in MM
	mmRs, err := u.muttermost.CreateUser(toMM(user))
	if err != nil {
		return nil, err
	}

	user.MMUserId = mmRs.Id
	user.Id = kit.NewId()

	// TODO: this shouldn't be here
	// create a private channel client-consultant for the client
	if user.Type == USER_TYPE_CLIENT {
		chRs, err := u.muttermost.CreateClientChannel(&mattermost.MMCreateClientChannelRequest{ClientUserId: user.MMUserId})
		if err != nil {
			return nil, err
		}
		user.MMChannelId = chRs.ChannelId
	}

	// save to storage
	dto, err := u.storage.CreateUser(toDto(user))
	if err != nil {
		return nil, err
	}

	user = fromDto(dto)

	return user, nil

}

func (u *userServiceImpl) GetByUsername(username string) *User {
	return fromDto(u.storage.GetByUsername(username))
}

func (u *userServiceImpl) GetByMMId(mmId string) *User {
	return fromDto(u.storage.GetByMMId(mmId))
}
