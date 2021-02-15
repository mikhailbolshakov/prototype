package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	pb "gitlab.medzdrav.ru/prototype/proto/chat"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"strings"
	"time"
)

type userServiceImpl struct {
	common.BaseService
	storage     domain.UserStorage
	chatService domain.ChatService
}

func NewUserService(storage domain.UserStorage, chatService domain.ChatService, queue queue.Queue) domain.UserService {

	s := &userServiceImpl{
		storage: storage,
		chatService: chatService,
	}
	s.BaseService = common.BaseService{Queue: queue}

	return s
}

func (u *userServiceImpl) newClient(ctx context.Context, user *domain.User) (*domain.User, error) {

	if user.ClientDetails == nil {
		return nil, fmt.Errorf("details isn't populated properly")
	}

	if user.ClientDetails.Phone == "" {
		return nil, fmt.Errorf("phone is empty")
	}

	var sex = map[string]bool{"M": true, "F": true}

	if _, ok := sex[user.ClientDetails.Sex]; !ok {
		return nil, fmt.Errorf("sex is incorrect")
	}

	if user.ClientDetails.FirstName == "" ||
		user.ClientDetails.LastName == "" ||
		user.ClientDetails.BirthDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) ||
		user.ClientDetails.BirthDate.After(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return nil, fmt.Errorf("pesonal data is incorrect")
	}

	if user.ClientDetails.PersonalAgreement == nil {
		user.ClientDetails.PersonalAgreement = &domain.PersonalAgreement{}
	}

	user.Username = user.ClientDetails.Phone

	if user.Groups == nil || len(user.Groups) == 0 {
		user.Groups = []string{domain.USER_GRP_CLIENT}
	}

	return user, nil

}

func (u *userServiceImpl) newConsultant(ctx context.Context, user *domain.User) (*domain.User, error) {

	if user.ConsultantDetails == nil {
		return nil, fmt.Errorf("details isn't populated properly")
	}

	if user.ConsultantDetails.Email == "" {
		return nil, fmt.Errorf("email is empty")
	}

	if user.ConsultantDetails.FirstName == "" ||
		user.ConsultantDetails.LastName == "" {
		return nil, fmt.Errorf("pesonal data is incorrect")
	}

	user.Username = strings.Split(user.ConsultantDetails.Email, "@")[0]

	if user.Groups == nil || len(user.Groups) == 0 {
		return nil, fmt.Errorf("groups aren't specified")
	}

	return user, nil

}

func (u *userServiceImpl) newExpert(ctx context.Context, user *domain.User) (*domain.User, error) {

	if user.ExpertDetails == nil {
		return nil, fmt.Errorf("details isn't populated properly")
	}

	if user.ExpertDetails.Email == "" {
		return nil, fmt.Errorf("email is empty")
	}

	if user.ExpertDetails.FirstName == "" ||
		user.ExpertDetails.LastName == "" {
		return nil, fmt.Errorf("pesonal data is incorrect")
	}

	user.Username = strings.Split(user.ExpertDetails.Email, "@")[0]

	if user.Groups == nil || len(user.Groups) == 0 {
		return nil, fmt.Errorf("groups aren't specified")
	}

	return user, nil

}

func (u *userServiceImpl) Create(ctx context.Context, user *domain.User) (*domain.User, error) {

	user.Id = kit.NewId()
	user.Status = domain.USER_STATUS_DRAFT

	var err error
	switch user.Type {
	case domain.USER_TYPE_CLIENT:
		user, err = u.newClient(ctx, user)
	case domain.USER_TYPE_CONSULTANT:
		user, err = u.newConsultant(ctx, user)
	case domain.USER_TYPE_EXPERT:
		user, err = u.newExpert(ctx, user)
	case domain.USER_TYPE_SUPERVISOR:
		return nil, errors.New("not implemented")
	default:
		return nil, fmt.Errorf("not supported user type %s", user.Type)
	}
	if err != nil {
		return nil, err
	}

	// check username uniqueness
	if usr := u.storage.GetByUsername(ctx, user.Username); usr != nil && usr.Id != "" {
		return nil, fmt.Errorf("username %s already exists", user.Username)
	}

	// save to storage
	user, err = u.storage.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// publish message
	if err := u.Publish(ctx, user, queue.QUEUE_TYPE_AT_LEAST_ONCE, "users.draft-created"); err != nil {
		return nil, err
	}

	return user, nil

}

func (u *userServiceImpl) GetByUsername(ctx context.Context, username string) *domain.User {
	return u.storage.GetByUsername(ctx, username)
}

func (u *userServiceImpl) GetByMMId(ctx context.Context, mmId string) *domain.User {
	return u.storage.GetByMMId(ctx, mmId)
}

func (u *userServiceImpl) Get(ctx context.Context, id string) *domain.User {
	return u.storage.Get(ctx, id)
}

func (u *userServiceImpl) Activate(ctx context.Context, userId string) (*domain.User, error) {

	user := u.storage.Get(ctx, userId)
	if user == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	return u.storage.UpdateStatus(ctx, userId, domain.USER_STATUS_ACTIVE, false)

}

func (u *userServiceImpl) Delete(ctx context.Context, userId string) (*domain.User, error) {

	user := u.storage.Get(ctx, userId)
	if user == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	return u.storage.UpdateStatus(ctx, userId, domain.USER_STATUS_DELETED, true)
}

func (u *userServiceImpl) SetClientDetails(ctx context.Context, userId string, details *domain.ClientDetails) (*domain.User, error) {

	user := u.storage.Get(ctx, userId)
	if user == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	if user.Type != domain.USER_TYPE_CLIENT {
		return nil, fmt.Errorf("user withid %s isn't a client", userId)
	}

	detB, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	return u.storage.UpdateDetails(ctx, userId, string(detB))

}

func (u *userServiceImpl) SetMMUserId(ctx context.Context, userId, mmId string) (*domain.User, error) {

	user := u.storage.Get(ctx, userId)
	if user == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	return u.storage.UpdateMMId(ctx, userId, mmId)

}

func (u *userServiceImpl) SetKKUserId(ctx context.Context, userId, kkId string) (*domain.User, error) {

	user := u.storage.Get(ctx, userId)
	if user == nil {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	return u.storage.UpdateKKId(ctx, userId, kkId)
}

func (u *userServiceImpl) Search(ctx context.Context, cr *domain.SearchCriteria) (*domain.SearchResponse, error) {

	if cr.PagingRequest == nil {
		cr.PagingRequest = &common.PagingRequest{}
	}

	if cr.Size == 0 {
		cr.Size = 100
	}

	if cr.Index == 0 {
		cr.Index = 1
	}

	response, err := u.storage.Search(ctx, cr)
	if err != nil {
		return nil, err
	}

	if len(response.Users) > 0 && cr.OnlineStatuses != nil && len(cr.OnlineStatuses) > 0 {
		var chatUserIds []string
		modifiedResponse := &domain.SearchResponse{Users: []*domain.User{}, PagingResponse: &common.PagingResponse{
			Total: response.Total,
			Index: response.Index,
		}}

		for _, u := range response.Users {
			chatUserIds = append(chatUserIds, u.MMUserId)
		}

		if mmStatuses, err := u.chatService.GetUsersStatuses(ctx, &pb.GetUsersStatusesRequest{ChatUserIds: chatUserIds}); err == nil {

			for _, user := range response.Users {

				for _, mmSt := range mmStatuses.Statuses {

					if mmSt.ChatUserId == user.MMUserId {

						for _, criteriaStatus := range cr.OnlineStatuses {

							if criteriaStatus == mmSt.Status {
								modifiedResponse.Users = append(modifiedResponse.Users, user)
							}

						}

					}

				}

			}

			if response.Total < cr.Size {
				modifiedResponse.Total = len(modifiedResponse.Users)
			}

			return modifiedResponse, nil

		} else {
			return nil, err
		}
	}

	return response, nil
}
