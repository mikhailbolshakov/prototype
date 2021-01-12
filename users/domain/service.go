package domain

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto/mm"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
	"log"
)

type UserService interface {
	Create(user *User) (*User, error)
	GetByUsername(username string) *User
	GetByMMId(mmId string) *User
}

type userServiceImpl struct {
	storage    storage.UserStorage
	mattermost mattermost.Service
	queue      queue.Queue
}

func NewUserService(storage storage.UserStorage, mmService mattermost.Service, queue queue.Queue) UserService {
	s := &userServiceImpl{
		storage:    storage,
		mattermost: mmService,
		queue:      queue,
	}

	return s
}

func (u *userServiceImpl) publish(user *User, topic string) {
	go func() {
		j, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = u.queue.Publish(topic, j)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()
}

func (u *userServiceImpl) Create(user *User) (*User, error) {

	// TODO: this shouldn't be here
	// send message to Queue, consumer in BP-service should create MM user and channel
	// user-service update mm_id and mm_channel_id on receiving an answer

	// create a new user in MM
	mmRs, err := u.mattermost.CreateUser(&mm.CreateUserRequest{
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return nil, err
	}

	user.MMUserId = mmRs.Id
	user.Id = kit.NewId()

	// create a private channel client-consultant for the client
	if user.Type == USER_TYPE_CLIENT {
		chRs, err := u.mattermost.CreateClientChannel(&mm.CreateClientChannelRequest{ClientUserId: user.MMUserId})
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
