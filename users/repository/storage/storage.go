package storage

import (
	"github.com/google/uuid"
	"gitlab.medzdrav.ru/prototype/kit"
	kitStorage "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"time"
)

type userGroup struct {
	kitStorage.BaseDto
	Id     string `gorm:"column:id"`
	UserId string `gorm:"column:user_id"`
	Type   string `gorm:"column:type"`
	Group  string `gorm:"column:group_code"`
}

type user struct {
	kitStorage.BaseDto
	Id       string       `gorm:"column:id"`
	Type     string       `gorm:"column:type"`
	Status   string       `gorm:"column:status"`
	Username string       `gorm:"column:username"`
	MMUserId string       `gorm:"column:mm_id"`
	KKUserId string       `gorm:"column:kk_id"`
	Details  string       `gorm:"column:details"`
	Groups   []*userGroup `gorm:"-"`
}

type storageImpl struct {
	c *container
}

func newStorage(c *container) *storageImpl {
	s := &storageImpl{c}
	return s
}

func (s *storageImpl) CreateUser(user *domain.User) (*domain.User, error) {

	dto := s.toUserDto(user)

	t := time.Now().UTC()
	dto.CreatedAt, dto.UpdatedAt = t, t

	result := s.c.Db.Instance.Create(dto)

	if result.Error != nil {
		return nil, result.Error
	}

	if err := s.addGroups(dto.Groups...); err != nil {
		return nil, err
	}

	return user, nil

}

func (s *storageImpl) updateField(userId, fieldName string, value interface{}) error {
	return s.c.Db.Instance.Model(&user{Id: userId}).
		Updates(map[string]interface{}{fieldName: value, "updated_at": time.Now().UTC()}).Error
}

func (s *storageImpl) UpdateStatus(userId, status string, isDeleted bool) (*domain.User, error) {

	var deletedAt *time.Time = nil
	if isDeleted {
		t := time.Now().UTC()
		deletedAt = &t
	}

	if err:=  s.c.Db.Instance.Model(&user{Id: userId}).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now().UTC(), "deleted_at": deletedAt}).Error; err != nil {

	}
	return s.Get(userId), nil
}

func (s *storageImpl) UpdateMMId(userId, mmId string) (*domain.User, error) {
	if err := s.updateField(userId, "mm_id", mmId); err != nil {
		return nil, err
	}
	return s.Get(userId), nil
}

func (s *storageImpl) UpdateKKId(userId, kkId string) (*domain.User, error) {
	if err := s.updateField(userId, "kk_id", kkId); err != nil {
		return nil, err
	}
	return s.Get(userId), nil
}

func (s *storageImpl) UpdateDetails(userId string, details string) (*domain.User, error) {
	if err := s.updateField(userId, "details", details); err != nil {
		return nil, err
	}
	return s.Get(userId), nil
}

func (s *storageImpl) getGroups(userId string) []*userGroup {
	var res []*userGroup
	if userId == "" {
		return res
	}
	s.c.Db.Instance.Where("user_id = ?::uuid", userId).Find(&res)
	return res
}

func (s *storageImpl) Get(id string) *domain.User {

	_, err := uuid.Parse(id)
	if err != nil {
		return s.GetByUsername(id)
	} else {
		dto := &user{Id: id}
		s.c.Db.Instance.First(dto)
		dto.Groups = s.getGroups(dto.Id)
		return s.toUserDomain(dto)
	}

}

func (s *storageImpl) GetByUsername(username string) *domain.User {
	dto := &user{}
	s.c.Db.Instance.Where("username = ? and deleted_at is null", username).First(&dto)
	dto.Groups = s.getGroups(dto.Id)
	return s.toUserDomain(dto)
}

func (s *storageImpl) GetByMMId(mmId string) *domain.User {
	dto := &user{}
	s.c.Db.Instance.Where("mm_id = ? and deleted_at is null", mmId).First(&dto)
	dto.Groups = s.getGroups(dto.Id)
	return s.toUserDomain(dto)
}


func (s *storageImpl) addGroups(groups ...*userGroup) error {

	t := time.Now().UTC()
	for _, g := range groups {
		g.CreatedAt, g.UpdatedAt = t, t
		if g.Id == "" {
			g.Id = kit.NewId()
		}
	}
	return s.c.Db.Instance.Create(groups).Error
}

func (s *storageImpl) RevokeGroups(groups ...*userGroup) error {
	t := time.Now().UTC()
	for _, g := range groups {
		g.DeletedAt = &t
	}
	return s.c.Db.Instance.Save(groups).Error
}
