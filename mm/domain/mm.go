package domain

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/tasks"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/queue_model"
)

type Service interface {
	GetUsersStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error)
	CreateUser(rq *CreateUserRequest) (*CreateUserResponse, error)
	CreateClientChannel(rq *CreateClientChannelRequest) (*CreateClientChannelResponse, error)
	GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error)
	SendTriggerPost(rq *SendTriggerPostRequest) error
	SubscribeUser(rq *SubscribeUserRequest) error
	TaskRemindMessageHandler(payload []byte) error
	MattermostPostMessageHandler(payload []byte) error
}

type serviceImpl struct {
	mmService    mattermost.Service
	usersService users.Service
	tasksService tasks.Service
	bpm          bpm.Engine
}

func NewService(mmService mattermost.Service, usersService users.Service, tasksService tasks.Service, bpm bpm.Engine) Service {

	s := &serviceImpl{
		mmService:    mmService,
		usersService: usersService,
		tasksService: tasksService,
		bpm:          bpm,
	}

	return s
}

func (s *serviceImpl) GetChannelsForUserAndMembers(rq *GetChannelsForUserAndMembersRequest) ([]string, error) {
	return s.mmService.GetChannelsForUserAndMembers(&mattermost.GetChannelsForUserAndMembersRequest{
		UserId:        rq.UserId,
		TeamName:      rq.TeamName,
		MemberUserIds: rq.MemberUserIds,
	})
}

func (s *serviceImpl) GetUsersStatuses(rq *GetUsersStatusesRequest) (*GetUsersStatusesResponse, error) {

	rs, err := s.mmService.GetUserStatuses(&mattermost.GetUsersStatusesRequest{UserIds: rq.UserIds})
	if err != nil {
		return nil, err
	}

	response := &GetUsersStatusesResponse{Statuses: []*UserStatus{}}
	for _, s := range rs.Statuses {
		response.Statuses = append(response.Statuses, &UserStatus{
			UserId: s.UserId,
			Status: s.Status,
		})
	}

	return response, nil

}

func (s *serviceImpl) CreateUser(rq *CreateUserRequest) (*CreateUserResponse, error) {

	rs, err := s.mmService.CreateUser(&mattermost.CreateUserRequest{
		TeamName: rq.TeamName,
		Username: rq.Username,
		Password: rq.Password,
		Email:    rq.Email,
	})

	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{Id: rs.Id}, nil
}

func (s *serviceImpl) CreateClientChannel(rq *CreateClientChannelRequest) (*CreateClientChannelResponse, error) {

	rs, err := s.mmService.CreateClientChannel(&mattermost.CreateClientChannelRequest{
		ClientUserId: rq.ClientUserId,
		TeamName:     rq.TeamName,
		DisplayName:  rq.DisplayName,
		Name:         rq.Name,
	})
	if err != nil {
		return nil, err
	}

	if rq.Subscribers != nil && len(rq.Subscribers) > 0 {
		for _, sbUserId := range rq.Subscribers {
			_, err = s.mmService.SubscribeUser(&mattermost.SubscribeUserRequest{
				UserId:    sbUserId,
				ChannelId: rs.ChannelId,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return &CreateClientChannelResponse{ChannelId: rs.ChannelId}, nil
}

func (s *serviceImpl) SubscribeUser(rq *SubscribeUserRequest) error {
	_, err := s.mmService.SubscribeUser(&mattermost.SubscribeUserRequest{
		UserId:    rq.UserId,
		ChannelId: rq.ChannelId,
	})
	return err
}

func (s *serviceImpl) MattermostPostMessageHandler(payload []byte) error{

	post := &mattermost.Post{}
	if err := json.Unmarshal(payload, post); err != nil {
		return err
	}

	// get user by MM user id
	user := s.usersService.GetByMMId(post.UserId)

	if user != nil && user.MMChannelId == post.ChannelId {

		variables := make(map[string]interface{})
		variables["userId"] = user.Id
		variables["username"] = user.Username
		variables["channelId"] = post.ChannelId
		variables["postTime"] = post.CreateAt

		_, err := s.bpm.StartProcess("p-client-request", variables)
		if err != nil {
			return err
		}

	}

	return nil

}

func (s *serviceImpl) TaskRemindMessageHandler(payload []byte) error {

	task := &queue_model.Task{}
	if err := json.Unmarshal(payload, task); err != nil {
		return err
	}

	reportedUser := s.usersService.GetByUsername(task.Reported.By)
	assigneeUser := s.usersService.GetByUsername(task.Assignee.User)

	dueDateStr := ""
	if task.DueDate != nil {
		dueDateStr = task.DueDate.Format("2006-01-02 15:04:05")
	}

	params := triggerPostParams{}
	params["task-num"] = task.Num
	params["due-date"] = dueDateStr
	if err := s.sendTriggerPost(TP_TASK_REMINDER, assigneeUser.MMId, task.ChannelId, params); err != nil {
		return err
	}

	if err := s.sendTriggerPost(TP_TASK_REMINDER, reportedUser.MMId, task.ChannelId, params); err != nil {
		return err
	}

	return nil
}
