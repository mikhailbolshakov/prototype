package storage

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/users/infrastructure"
	"math"
	"time"
)

type UserStorage interface {
	CreateUser(u *User) (*User, error)
	GetByUsername(username string) *User
	GetByMMId(mmId string) *User
	Get(id string) *User
	Search(cr *SearchCriteria) (*SearchResponse, error)
	UpdateStatus(userId, status string, isDeleted bool) (*User, error)
	UpdateDetails(userId string, details string) (*User, error)
	UpdateMMId(userId, mmId string) (*User, error)
	UpdateKKId(userId, kkId string) (*User, error)
}

type storageImpl struct {
	infr *infrastructure.Container
}

func NewStorage(infr *infrastructure.Container) UserStorage {
	s := &storageImpl{
		infr: infr,
	}
	return s
}

func (s *storageImpl) CreateUser(user *User) (*User, error) {

	t := time.Now().UTC()
	user.CreatedAt, user.UpdatedAt = t, t

	result := s.infr.Db.Instance.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s *storageImpl) updateField(userId, fieldName string, value interface{}) error {
	return s.infr.Db.Instance.Model(&User{Id: userId}).
		Updates(map[string]interface{}{fieldName: value, "updated_at": time.Now().UTC()}).Error
}

func (s *storageImpl) UpdateStatus(userId, status string, isDeleted bool) (*User, error) {

	var deletedAt *time.Time = nil
	if isDeleted {
		t := time.Now().UTC()
		deletedAt = &t
	}

	if err:=  s.infr.Db.Instance.Model(&User{Id: userId}).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now().UTC(), "deleted_at": deletedAt}).Error; err != nil {

	}
	return s.Get(userId), nil
}

func (s *storageImpl) UpdateMMId(userId, mmId string) (*User, error) {
	if err := s.updateField(userId, "mm_id", mmId); err != nil {
		return nil, err
	}
	return s.Get(userId), nil
}

func (s *storageImpl) UpdateKKId(userId, kkId string) (*User, error) {
	if err := s.updateField(userId, "kk_id", kkId); err != nil {
		return nil, err
	}
	return s.Get(userId), nil
}

func (s *storageImpl) UpdateDetails(userId string, details string) (*User, error) {
	if err := s.updateField(userId, "details", details); err != nil {
		return nil, err
	}
	return s.Get(userId), nil
}

func (s *storageImpl) Get(id string) *User{

	_, err := uuid.Parse(id)
	if err != nil {
		return s.GetByUsername(id)
	} else {
		user := &User{Id: id}
		s.infr.Db.Instance.First(user)
		return user
	}

}

func (s *storageImpl) GetByUsername(username string) *User {
	user := &User{}
	s.infr.Db.Instance.Where("username = ?", username).First(&user)
	return user
}

func (s *storageImpl) GetByMMId(mmId string) *User {
	user := &User{}
	s.infr.Db.Instance.Where("mm_id = ?", mmId).First(&user)
	return user
}

func (s *storageImpl) Search(cr *SearchCriteria) (*SearchResponse, error) {

	response := &SearchResponse{
		PagingResponse: &common.PagingResponse{
			Total: 0,
			Index: 0,
		},
		Users: []*User{},
	}

	selectClause := `*`

	query := s.infr.Db.Instance.
		Table(`users u`).
		Where(`u.deleted_at is null`)

	if cr.Username != "" {
		query = query.Where(`u.username = ?`, cr.Username)
	}

	if cr.Status != "" {
		query = query.Where(`u.status = ?`, cr.Status)
	}

	if cr.MMChannelId != "" {
		query = query.Where(`(u.details -> 'mmChannelId')::varchar = ?`, fmt.Sprintf(`"%s"`, cr.MMChannelId))
	}

	if cr.MMId != "" {
		query = query.Where(`u.mm_id = ?`, cr.MMId)
	}

	if cr.Email != "" {
		query = query.Where(`(u.details -> 'email')::varchar = ?`, fmt.Sprintf(`"%s"`, cr.Email))
	}

	if cr.Phone != "" {
		query = query.Where(`(u.details -> 'phone')::varchar = ?`, fmt.Sprintf(`"%s"`, cr.Phone))
	}

	if cr.UserType != "" {
		query = query.Where(`u.type = ?`, cr.UserType)
	}

	// paging
	var totalCount int64
	var offset int

	query.Count(&totalCount)

	if totalCount > int64(cr.Size) {
		offset = (cr.Index - 1) * cr.Size
	}

	response.PagingResponse.Total = int(math.Ceil(float64(totalCount) / float64(cr.Size)))
	response.PagingResponse.Index = cr.Index

	query = query.Select(selectClause).Offset(offset).Limit(cr.Size)

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &User{}
		_ = s.infr.Db.Instance.ScanRows(rows, task)
		response.Users = append(response.Users, task)
	}

	return response, nil
}
