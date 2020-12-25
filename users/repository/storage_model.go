package repository

import kit "gitlab.medzdrav.ru/prototype/kit/storage"

type User struct {
	kit.BaseDto
	Id        string `gorm:"column:id"`
	Type      string `gorm:"column:type"`
	Username  string `gorm:"column:username"`
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Phone     string `gorm:"column:phone"`
	Email     string `gorm:"column:email"`
	MMUserId  string `gorm:"column:mm_id"`
}
