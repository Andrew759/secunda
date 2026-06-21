package user

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Phone     string    `json:"phone" gorm:"type:varchar(30)"`
	Name      string    `json:"name" gorm:"type:varchar(256);not null"`
	Surname   string    `json:"surname" gorm:"type:varchar(256);not null"`
	Login     string    `json:"login" gorm:"type:varchar(256);unique; not null"`
	Password  *string   `json:"password" gorm:"type:varchar(1024)"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"type:timestamp;not null"`
}

type AlreadyExistErr struct {
	Field string
	Value string
}

func (e *AlreadyExistErr) Error() string {
	return fmt.Sprintf("user with %s '%s' already exists", e.Field, e.Value)
}

var NotFoundErr = errors.New("user not found")

func (u User) TableName() string {
	return "users"
}

func CreateUser(db *gorm.DB, u *User) error {
	var existing User
	err := db.Where("login = ? OR phone = ?", u.Login, u.Phone).First(&existing).Error

	if err == nil {
		if existing.Login == u.Login {
			return &AlreadyExistErr{Field: "login", Value: u.Login}
		}
		if existing.Phone == u.Phone {
			return &AlreadyExistErr{Field: "phone", Value: u.Phone}
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFoundErr
	}

	return db.Create(u).Error
}

func GetUserById(db *gorm.DB, id int) (User, error) {
	var user User
	result := db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, NotFoundErr
	}

	return user, result.Error
}
